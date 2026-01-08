# Example: External Custom Fields Management
#
# This pattern is used when custom fields are managed outside of Terraform
# (e.g., via NetBox UI, NetBox API, or external automation tools).
# By omitting the custom_fields attribute, Terraform will preserve ALL
# existing custom fields in NetBox during updates.

resource "netbox_site" "external" {
  name   = "External Management Site"
  slug   = "external-site"
  status = "active"
}

resource "netbox_manufacturer" "external" {
  name = "External Manufacturer"
  slug = "external-manufacturer"
}

resource "netbox_device_type" "external" {
  model        = "External Model"
  slug         = "external-model"
  manufacturer = netbox_manufacturer.external.id
  u_height     = 1
}

resource "netbox_device_role" "external" {
  name  = "External Role"
  slug  = "external-role"
  color = "00ff00"
}

resource "netbox_device" "external_cf_management" {
  name        = "external-cf-device"
  device_type = netbox_device_type.external.model
  role        = netbox_device_role.external.name
  site        = netbox_site.external.name
  status      = "active"

  # custom_fields attribute is intentionally omitted
  # All custom fields set in NetBox are preserved during Terraform operations
  # This includes:
  # - Updates to other device attributes (name, status, etc.)
  # - Import operations
  # - Refresh operations

  tags = [
    "external-custom-fields"
  ]
}
