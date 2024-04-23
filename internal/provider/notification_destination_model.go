package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type notificationDestination struct {
	ID       types.String                       `tfsdk:"id"`
	Name     types.String                       `tfsdk:"name"`
	Resource notificationDestinationAWSResource `tfsdk:"resource"`
}

type notificationDestinationAWSResource struct {
	SQS         *notificationDestinationAWSResourceSQS         `tfsdk:"sqs"`
	EventBridge *notificationDestinationAWSResourceEventBridge `tfsdk:"event_bridge"`
}

type notificationDestinationAWSResourceSQS struct {
	ARN types.String `tfsdk:"arn"`
}

type notificationDestinationAWSResourceEventBridge struct {
	Name      types.String `tfsdk:"name"`
	Region    types.String `tfsdk:"region"`
	AccountID types.String `tfsdk:"account_id"`
}

type notificationDestinationResourceModel struct {
	Region types.String `tfsdk:"region"`
	notificationDestination
}

type notificationDestinationsDataSourceModel struct {
	Region       types.String              `tfsdk:"region"`
	Destinations []notificationDestination `tfsdk:"destinations"`
}
