# Event Rule Data Source Test Configuration
# This file tests the netbox_event_rule data source

terraform {
  required_version = ">= 1.0"
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {}

# Create a webhook for the event rule to use
resource "netbox_webhook" "test_ds" {
  name        = "test-webhook-event-rule-ds"
  payload_url = "https://example.com/webhook/ds"
}

# Create an event rule for the data source to look up
resource "netbox_event_rule" "source" {
  name               = "test-event-rule-ds"
  object_types       = ["dcim.device"]
  event_types        = ["object_created", "object_updated"]
  action_type        = "webhook"
  action_object_type = "extras.webhook"
  action_object_id   = netbox_webhook.test_ds.id
  enabled            = true
  description        = "Test event rule for data source lookup"
}

# Look up the event rule by ID
data "netbox_event_rule" "test" {
  id = netbox_event_rule.source.id
}
