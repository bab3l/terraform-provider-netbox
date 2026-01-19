resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_type" "test" {
  model        = "Test Model"
  slug         = "test-model"
  manufacturer = netbox_manufacturer.test.id
  u_height     = 1
}

resource "netbox_device_role" "test" {
  name  = "Test Role"
  slug  = "test-role"
  color = "ff0000"
}

resource "netbox_device" "test" {
  name        = "test-device-1"
  device_type = netbox_device_type.test.model
  role        = netbox_device_role.test.name
  site        = netbox_site.test.name
  status      = "active"

  serial    = "1234567890"
  asset_tag = "asset-123"

  # Partial custom fields management (recommended pattern)
  # Only the custom fields specified here are managed by Terraform
  # Other custom fields set in NetBox (via UI, API, or automation) are preserved
  custom_fields = [
    {
      name  = "environment"
      value = "production"
    },
    {
      name  = "owner_team"
      value = "network-ops"
    }
  ]

  tags = [
    "managed-by-terraform",
    "production"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_device.test
  id = "123"

  identity = {
    custom_fields = [
      "environment:text",
      "owner_team:text",
    ]
  }
}
