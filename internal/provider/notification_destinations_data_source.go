package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/brandedtech/sp-api-sdk/notifications"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &notificationDestinationsDatasource{}
	_ datasource.DataSourceWithConfigure = &notificationDestinationsDatasource{}
)

func NewNotificationDestinationsDatasource() datasource.DataSource {
	return &notificationDestinationsDatasource{}
}

type notificationDestinationsDatasource struct {
	providerData *SPAPIProviderData
}

func (d *notificationDestinationsDatasource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*SPAPIProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *sp.SellingPartner, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.providerData = providerData
}

func (d *notificationDestinationsDatasource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_destinations"
}

func (d *notificationDestinationsDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Computed: true,
			},
			"destinations": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"resource": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"sqs": schema.SingleNestedAttribute{
									Computed: true,
									Optional: true,
									Attributes: map[string]schema.Attribute{
										"arn": schema.StringAttribute{
											Computed: true,
										},
									},
								},
								"event_bridge": schema.SingleNestedAttribute{
									Computed: true,
									Optional: true,
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Computed: true,
										},
										"region": schema.StringAttribute{
											Computed: true,
										},
										"account_id": schema.StringAttribute{
											Computed: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (n *notificationDestinationsDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state notificationDestinationsDataSourceModel

	client, err := notifications.NewClientWithResponses(RegionToEndpoint(state.Region.ValueString()),
		notifications.WithRequestBefore(func(ctx context.Context, req *http.Request) error {
			return n.providerData.Grantless.AuthorizeRequestWithScope(req, "sellingpartnerapi::notifications")
		}),
	)

	if err != nil {
		resp.Diagnostics.AddError("Error creating notifications client", err.Error())
		return
	}

	destinations, err := client.GetDestinationsWithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error getting destinations", err.Error())
		return
	}

	for _, destination := range *destinations.Model.Payload {
		model := notificationDestination{
			ID:   types.StringValue(destination.DestinationId),
			Name: types.StringValue(destination.Name),
		}

		if destination.Resource.Sqs != nil {
			model.Resource.SQS = &notificationDestinationAWSResourceSQS{
				ARN: types.StringValue(destination.Resource.Sqs.Arn),
			}
		}

		if destination.Resource.EventBridge != nil {
			model.Resource.EventBridge = &notificationDestinationAWSResourceEventBridge{
				Name:      types.StringValue(destination.Resource.EventBridge.Name),
				Region:    types.StringValue(destination.Resource.EventBridge.Region),
				AccountID: types.StringValue(destination.Resource.EventBridge.AccountId),
			}
		}

		state.Destinations = append(state.Destinations, model)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
