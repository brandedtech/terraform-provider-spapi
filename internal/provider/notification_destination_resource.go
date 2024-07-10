package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/brandedtech/sp-api-sdk/notifications"
	sp "github.com/brandedtech/sp-api-sdk/pkg/selling-partner"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &notificationDestinationResource{}
	_ resource.ResourceWithConfigure = &notificationDestinationResource{}
)

// NewNotificationDestinationResource is a helper function to simplify the provider implementation.
func NewNotificationDestinationResource() resource.Resource {
	return &notificationDestinationResource{}
}

// orderResource is the resource implementation.
type notificationDestinationResource struct {
	sellingPartner *sp.SellingPartner
}

// Metadata returns the resource type name.
func (r *notificationDestinationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_destination"
}

func (r *notificationDestinationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	sellingPartner, ok := req.ProviderData.(*sp.SellingPartner)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *sp.SellingPartner, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.sellingPartner = sellingPartner
}

// Schema defines the schema for the resource.
func (r *notificationDestinationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"resource": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"sqs": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"arn": schema.StringAttribute{
								Required: true,
							},
						},
					},
					"event_bridge": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Computed: true,
							},
							"region": schema.StringAttribute{
								Required: true,
							},
							"account_id": schema.StringAttribute{
								Required: true,
							},
						},
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *notificationDestinationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan notificationDestinationModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := notifications.NewClientWithResponses("https://sellingpartnerapi-na.amazon.com",
		notifications.WithRequestBefore(func(ctx context.Context, req *http.Request) error {
			return r.sellingPartner.AuthorizeRequestWithScope(req, "sellingpartnerapi::notifications")
		}),
	)

	if err != nil {
		resp.Diagnostics.AddError("Error creating notifications client", err.Error())
		return
	}

	body := notifications.CreateDestinationJSONRequestBody{
		Name: plan.Name.ValueString(),
	}

	if plan.Resource.SQS != nil {
		body.ResourceSpecification = notifications.DestinationResourceSpecification{
			Sqs: &notifications.SqsResource{
				Arn: plan.Resource.SQS.ARN.ValueString(),
			},
		}
	} else if plan.Resource.EventBridge != nil {
		body.ResourceSpecification = notifications.DestinationResourceSpecification{
			EventBridge: &notifications.EventBridgeResourceSpecification{
				Region:    plan.Resource.EventBridge.Region.ValueString(),
				AccountId: plan.Resource.EventBridge.AccountID.ValueString(),
			},
		}
	}

	destination, err := client.CreateDestinationWithResponse(ctx, body)

	if err != nil {
		resp.Diagnostics.AddError("Error creating destination", err.Error())
		return
	}

	plan.ID = types.StringValue(destination.Model.Payload.DestinationId)

	if plan.Resource.EventBridge != nil {
		plan.Resource.EventBridge.Name = types.StringValue(destination.Model.Payload.Resource.EventBridge.Name)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *notificationDestinationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notificationDestinationModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := notifications.NewClientWithResponses("https://sellingpartnerapi-na.amazon.com",
		notifications.WithRequestBefore(func(ctx context.Context, req *http.Request) error {
			return r.sellingPartner.AuthorizeRequestWithScope(req, "sellingpartnerapi::notifications")
		}),
	)

	if err != nil {
		resp.Diagnostics.AddError("Error creating notifications client", err.Error())
		return
	}

	destination, err := client.GetDestinationWithResponse(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting destination", err.Error())
		return
	}

	state.Name = types.StringValue(destination.Model.Payload.Name)

	if destination.Model.Payload.Resource.Sqs != nil {
		state.Resource.SQS.ARN = types.StringValue(destination.Model.Payload.Resource.Sqs.Arn)
	} else if destination.Model.Payload.Resource.EventBridge != nil {
		state.Resource.EventBridge.Region = types.StringValue(destination.Model.Payload.Resource.EventBridge.Region)
		state.Resource.EventBridge.AccountID = types.StringValue(destination.Model.Payload.Resource.EventBridge.AccountId)
		state.Resource.EventBridge.Name = types.StringValue(destination.Model.Payload.Resource.EventBridge.Name)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *notificationDestinationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *notificationDestinationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state notificationDestinationModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := notifications.NewClientWithResponses("https://sellingpartnerapi-na.amazon.com",
		notifications.WithRequestBefore(func(ctx context.Context, req *http.Request) error {
			return r.sellingPartner.AuthorizeRequestWithScope(req, "sellingpartnerapi::notifications")
		}),
	)

	if err != nil {
		resp.Diagnostics.AddError("Error creating notifications client", err.Error())
		return
	}

	destinationDeleteResp, err := client.DeleteDestinationWithResponse(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting destination", err.Error())
		return
	}

	if destinationDeleteResp.Model.Errors != nil {
		for _, error := range *destinationDeleteResp.Model.Errors {
			resp.Diagnostics.AddError("Error deleting destination", error.Message)
		}
		return
	}
}
