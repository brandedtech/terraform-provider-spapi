package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type notificationSubscription struct {
	ID               types.String `tfsdk:"id"`
	NotificationType types.String `tfsdk:"notification_type"`
	PayloadVersion   types.String `tfsdk:"payload_version"`
	DestinationID    types.String `tfsdk:"destination_id"`
}

type notificationSubscriptionResourceModel struct {
	Region types.String `tfsdk:"region"`
	notificationSubscription
}
