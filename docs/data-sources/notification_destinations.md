---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "spapi_notification_destinations Data Source - spapi"
subcategory: ""
description: |-
  
---

# spapi_notification_destinations (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `destinations` (Attributes List) (see [below for nested schema](#nestedatt--destinations))

<a id="nestedatt--destinations"></a>
### Nested Schema for `destinations`

Read-Only:

- `id` (String)
- `name` (String)
- `resource` (Attributes) (see [below for nested schema](#nestedatt--destinations--resource))

<a id="nestedatt--destinations--resource"></a>
### Nested Schema for `destinations.resource`

Optional:

- `event_bridge` (Attributes) (see [below for nested schema](#nestedatt--destinations--resource--event_bridge))
- `sqs` (Attributes) (see [below for nested schema](#nestedatt--destinations--resource--sqs))

<a id="nestedatt--destinations--resource--event_bridge"></a>
### Nested Schema for `destinations.resource.event_bridge`

Read-Only:

- `account_id` (String)
- `name` (String)
- `region` (String)


<a id="nestedatt--destinations--resource--sqs"></a>
### Nested Schema for `destinations.resource.sqs`

Read-Only:

- `arn` (String)