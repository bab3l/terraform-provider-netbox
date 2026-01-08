# Example: Explicit Custom Field Removal
#
# This pattern demonstrates how to explicitly remove specific custom fields
# while preserving others. Setting a custom field's value to an empty string
# removes it from NetBox, but other custom fields not listed remain intact.

resource "netbox_site" "removal" {
  name   = "Removal Example Site"
  slug   = "removal-site"
  status = "active"
}

resource "netbox_manufacturer" "removal" {
  name = "Removal Manufacturer"
  slug = "removal-manufacturer"
}

resource "netbox_device_type" "removal" {
  model        = "Removal Model"
  slug         = "removal-model"
  manufacturer = netbox_manufacturer.removal.id
  u_height     = 1
}

resource "netbox_device_role" "removal" {
  name  = "Removal Role"
  slug  = "removal-role"
  color = "0000ff"
}

resource "netbox_device" "explicit_removal" {
  name        = "removal-device"
  device_type = netbox_device_type.removal.model
  role        = netbox_device_role.removal.name
  site        = netbox_site.removal.name
  status      = "active"

  # To remove a specific custom field, set its value to an empty string
  # In this example:
  # - "environment" field is removed from NetBox
  # - "cost_center" field is still managed with a value
  # - Any other custom fields in NetBox are preserved (not listed here)
  custom_fields = [
    {
      name  = "environment"
      value = "" # Empty value removes this field from NetBox
    },
    {
      name  = "cost_center"
      value = "CC-12345"
    }
  ]

  tags = [
    "explicit-removal-example"
  ]
}
