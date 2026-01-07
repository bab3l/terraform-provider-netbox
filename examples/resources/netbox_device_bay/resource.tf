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
  manufacturer = netbox_manufacturer.test.name
  u_height     = 1
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device" "test" {
  name        = "test-device-1"
  device_type = netbox_device_type.test.model
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

resource "netbox_device_bay" "test" {
  name        = "Device Bay 1"
  device      = netbox_device.test.name
  description = "Expansion slot for modular device"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "bay_type"
      value = "network-module"
    },
    {
      name  = "max_power_watts"
      value = "50"
    }
  ]

  tags = [
    "device-bay",
    "expansion-slot"
  ]
}
