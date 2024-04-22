package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type notificationDestinationModel struct {
	ID       types.String                         `tfsdk:"id"`
	Name     types.String                         `tfsdk:"name"`
	Resource notificationDestinationResourceModel `tfsdk:"resource"`
}

type notificationDestinationResourceModel struct {
	SQS         *notificationDestinationResourceSQSModel         `tfsdk:"sqs"`
	EventBridge *notificationDestinationResourceEventBridgeModel `tfsdk:"event_bridge"`
}

type notificationDestinationResourceSQSModel struct {
	ARN types.String `tfsdk:"arn"`
}

type notificationDestinationResourceEventBridgeModel struct {
	Name      types.String `tfsdk:"name"`
	Region    types.String `tfsdk:"region"`
	AccountID types.String `tfsdk:"account_id"`
}
