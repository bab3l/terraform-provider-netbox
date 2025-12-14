# Notification Group Data Source Test Configuration
# This file tests the netbox_notification_group data source

terraform {
  required_version = ">= 1.0"
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {}

# Create a notification group for the data source to look up
resource "netbox_notification_group" "source" {
  name        = "test-notification-group-ds"
  description = "Test notification group for data source lookup"
}

# Look up the notification group by ID
data "netbox_notification_group" "test" {
  id = netbox_notification_group.source.id
}
