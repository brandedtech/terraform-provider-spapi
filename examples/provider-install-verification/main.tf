terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.45.0"
    }
    spapi = {
      source = "brandedtech/spapi"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

provider "spapi" {}

resource "spapi_notification_destination" "event_bridge" {
  name = var.spapi_account_id
  resource = {
    event_bridge = {
      region     = "us-east-1"
      account_id = "574253385567"
    }
  }
}

data "aws_cloudwatch_event_source" "event_bridge" {
  name_prefix = spapi_notification_destination.event_bridge.name
}

resource "spapi_notification_subscription" "event_bridge_BRANDED_ITEM_CONTENT_CHANGE" {
  notification_type = "BRANDED_ITEM_CONTENT_CHANGE"
  destination_id    = spapi_notification_destination.event_bridge.id
}
