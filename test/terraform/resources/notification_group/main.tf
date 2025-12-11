# Notification Group Resource Test Configuration
# This file tests the netbox_notification_group resource

terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {}

# Basic notification group - minimal configuration
resource "netbox_notification_group" "basic" {
  name = "test-notification-group-basic"
}

# Notification group with description
resource "netbox_notification_group" "with_description" {
  name        = "test-notification-group-desc"
  description = "Test notification group with description"
}
