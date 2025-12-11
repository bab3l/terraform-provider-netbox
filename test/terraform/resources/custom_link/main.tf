terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

# Basic custom link
resource "netbox_custom_link" "device_docs" {
  name         = "device_docs"
  object_types = ["dcim.device"]
  link_text    = "Documentation"
  link_url     = "https://docs.example.com/devices/{{ object.name }}"
}

# Custom link with full options
resource "netbox_custom_link" "full_example" {
  name         = "full_link"
  object_types = ["dcim.device", "dcim.site"]
  enabled      = true
  link_text    = "View in External System"
  link_url     = "https://external.example.com/{{ object.name }}/{{ object.id }}"
  weight       = 100
  group_name   = "External Links"
  button_class = "blue"
  new_window   = true
}

# Custom link for sites with conditional display
resource "netbox_custom_link" "site_map" {
  name         = "site_map"
  object_types = ["dcim.site"]
  link_text    = "{% if object.latitude %}View on Map{% else %}No Location{% endif %}"
  link_url     = "https://maps.google.com/?q={{ object.latitude }},{{ object.longitude }}"
  button_class = "green"
}

# Custom link for multiple object types
resource "netbox_custom_link" "asset_tracker" {
  name         = "asset_tracker"
  object_types = ["dcim.device", "virtualization.virtualmachine"]
  link_text    = "Track Asset: {{ object.name }}"
  link_url     = "https://assets.example.com/track?id={{ object.serial }}"
  weight       = 50
  new_window   = true
}

output "device_docs_id" {
  value = netbox_custom_link.device_docs.id
}

output "full_example_id" {
  value = netbox_custom_link.full_example.id
}

output "site_map_id" {
  value = netbox_custom_link.site_map.id
}

output "asset_tracker_id" {
  value = netbox_custom_link.asset_tracker.id
}
