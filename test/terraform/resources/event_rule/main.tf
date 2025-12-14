# Event Rule Resource Test Configuration
# This file tests the netbox_event_rule resource

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
resource "netbox_webhook" "test_basic" {
  name        = "test-webhook-event-rule-basic"
  payload_url = "https://example.com/webhook/basic"
}

resource "netbox_webhook" "test_complete" {
  name        = "test-webhook-event-rule-complete"
  payload_url = "https://example.com/webhook/complete"
}

# Create tags for the complete event rule
resource "netbox_tag" "event_rule_test" {
  name  = "event-rule-test"
  slug  = "event-rule-test"
  color = "ff5722"
}

# Basic event rule - minimal configuration
resource "netbox_event_rule" "basic" {
  name               = "test-event-rule-basic"
  object_types       = ["dcim.device"]
  event_types        = ["object_created"]
  action_type        = "webhook"
  action_object_type = "extras.webhook"
  action_object_id   = netbox_webhook.test_basic.id
}

# Event rule with multiple object types and event types
resource "netbox_event_rule" "multiple_types" {
  name               = "test-event-rule-multiple"
  object_types       = ["dcim.device", "dcim.site", "ipam.ipaddress"]
  event_types        = ["object_created", "object_updated", "object_deleted"]
  action_type        = "webhook"
  action_object_type = "extras.webhook"
  action_object_id   = netbox_webhook.test_basic.id
  enabled            = true
}

# Complete event rule with all options
resource "netbox_event_rule" "complete" {
  name               = "test-event-rule-complete"
  object_types       = ["dcim.device", "dcim.site"]
  event_types        = ["object_created", "object_updated"]
  action_type        = "webhook"
  action_object_type = "extras.webhook"
  action_object_id   = netbox_webhook.test_complete.id
  enabled            = true
  description        = "Complete event rule with all configuration options"

  tags = [
    {
      name  = netbox_tag.event_rule_test.name
      slug  = netbox_tag.event_rule_test.slug
      color = netbox_tag.event_rule_test.color
    }
  ]
}

# Disabled event rule
resource "netbox_event_rule" "disabled" {
  name               = "test-event-rule-disabled"
  object_types       = ["dcim.device"]
  event_types        = ["object_created"]
  action_type        = "webhook"
  action_object_type = "extras.webhook"
  action_object_id   = netbox_webhook.test_basic.id
  enabled            = false
  description        = "This event rule is disabled"
}
