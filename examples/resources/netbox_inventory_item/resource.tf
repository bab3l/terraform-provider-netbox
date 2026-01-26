resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_device_role" "test" {
  name  = "Test Role"
  slug  = "test-role"
  color = "ff0000"
}

resource "netbox_device_type" "test" {
  model        = "Test Model"
  slug         = "test-model"
  manufacturer = netbox_manufacturer.test.id
  u_height     = 1
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device" "test" {
  name        = "test-device-1"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

resource "netbox_inventory_item" "test" {
  name        = "Inventory Item 1"
  device      = netbox_device.test.id
  description = "Power supply module"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "serial_number"
      value = "PSU-SN-12345"
    },
    {
      name  = "purchase_date"
      value = "2024-01-15"
    },
    {
      name  = "warranty_expiry"
      value = "2027-01-15"
    },
    {
      name  = "asset_tag"
      value = "INV-PSU-001"
    }
  ]

  tags = [
    "power-supply",
    "warranty-tracked"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_inventory_item.test
  id = "123"

  identity = {
    custom_fields = [
      "serial_number:text",
      "purchase_date:date",
      "warranty_expiry:date",
      "asset_tag:text",
    ]
  }
}
