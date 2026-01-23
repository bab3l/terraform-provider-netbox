# Create a basic tag
resource "netbox_tag" "production" {
  name = "Production"
  slug = "production"
}

# Create a tag with a custom color
resource "netbox_tag" "critical" {
  name        = "Critical"
  slug        = "critical"
  color       = "ff0000"
  description = "Critical infrastructure requiring special attention"
}

# Create a tag restricted to specific object types
resource "netbox_tag" "network_device" {
  name         = "Network Device"
  slug         = "network-device"
  color        = "2196f3"
  description  = "Tag for network infrastructure devices"
  object_types = ["dcim.device", "dcim.interface"]
}

# Create an environment tag
resource "netbox_tag" "staging" {
  name        = "Staging"
  slug        = "staging"
  color       = "ff9800"
  description = "Staging environment resources"
}

import {
  to = netbox_tag.production
  id = "123"
}
