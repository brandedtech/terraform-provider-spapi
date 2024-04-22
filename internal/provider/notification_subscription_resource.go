package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/brandedtech/sp-api-sdk/notifications"
	sp "github.com/brandedtech/sp-api-sdk/pkg/selling-partner"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &notificationSubscriptionResource{}
	_ resource.ResourceWithConfigure = &notificationSubscriptionResource{}
)

// NewNotificationSubscriptionResource is a helper function to simplify the provider implementation.
func NewNotificationSubscriptionResource() resource.Resource {
	return &notificationSubscriptionResource{}
}

// orderResource is the resource implementation.
type notificationSubscriptionResource struct {
	sellingPartner *sp.SellingPartner
}

// Metadata returns the resource type name.
func (r *notificationSubscriptionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_subscription"
}

func (r *notificationSubscriptionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *notificationSubscriptionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"notification_type": schema.StringAttribute{
				Required: true,
			},
			"payload_version": schema.StringAttribute{
				Computed: true,
				Default:  stringdefault.StaticString("1.0"),
			},
			"destination_id": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *notificationSubscriptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan notificationSubscriptionModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := notifications.NewClientWithResponses("https://sellingpartnerapi-na.amazon.com",
		notifications.WithRequestBefore(func(ctx context.Context, req *http.Request) error {
			return r.sellingPartner.AuthorizeRequest(req)
		}),
	)

	if err != nil {
		resp.Diagnostics.AddError("Error creating notifications client", err.Error())
		return
	}

	payloadVersion := plan.PayloadVersion.ValueString()
	destinationId := plan.DestinationID.ValueString()

	body := notifications.CreateSubscriptionJSONRequestBody{
		PayloadVersion: &payloadVersion,
		DestinationId:  &destinationId,
	}

	subscription, err := client.CreateSubscriptionWithResponse(ctx, notifications.NotificationType(plan.NotificationType.ValueString()), body)

	if err != nil {
		resp.Diagnostics.AddError("Error creating subscription", err.Error())
		return
	}

	plan.ID = types.StringValue(subscription.Model.Payload.SubscriptionId)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *notificationSubscriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notificationSubscriptionModel

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

	subscription, err := client.GetSubscriptionByIdWithResponse(ctx, notifications.NotificationType(state.NotificationType.ValueString()), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting subscription", err.Error())
		return
	}

	state.PayloadVersion = types.StringValue(subscription.Model.Payload.PayloadVersion)
	state.DestinationID = types.StringValue(subscription.Model.Payload.DestinationId)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *notificationSubscriptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *notificationSubscriptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
