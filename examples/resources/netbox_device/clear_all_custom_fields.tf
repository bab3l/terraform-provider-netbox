# Example: Complete Custom Fields Removal
#
# This pattern demonstrates how to remove ALL custom fields from a resource.
# Setting custom_fields to an empty list will clear all custom fields in NetBox.

resource "netbox_site" "clear" {
  name   = "Clear All Site"
  slug   = "clear-site"
  status = "active"
}

resource "netbox_manufacturer" "clear" {
  name = "Clear Manufacturer"
  slug = "clear-manufacturer"
}

resource "netbox_device_type" "clear" {
  model        = "Clear Model"
  slug         = "clear-model"
  manufacturer = netbox_manufacturer.clear.id
  u_height     = 1
}

resource "netbox_device_role" "clear" {
  name  = "Clear Role"
  slug  = "clear-role"
  color = "ffff00"
}

resource "netbox_device" "clear_all_custom_fields" {
  name        = "clear-all-device"
  device_type = netbox_device_type.clear.model
  role        = netbox_device_role.clear.name
  site        = netbox_site.clear.name
  status      = "active"

  # Empty list explicitly clears ALL custom fields from NetBox
  # Use this when you want to ensure no custom fields exist on the resource
  # Note: This is different from omitting custom_fields entirely, which
  # preserves all existing custom fields
  custom_fields = []

  tags = [
    "no-custom-fields"
  ]
}
