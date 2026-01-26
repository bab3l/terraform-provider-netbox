resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_device_role" "test" {
  name    = "Test Role"
  slug    = "test-role"
  color   = "ff0000"
  vm_role = false
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

resource "netbox_device" "test" {
  name        = "Test Device"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_service" "test" {
  name        = "SSH"
  protocol    = "tcp"
  ports       = [22]
  device      = netbox_device.test.id
  description = "SSH service for remote administration"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "service_owner"
      value = "infrastructure-team"
    },
    {
      name  = "monitoring_enabled"
      value = "true"
    },
    {
      name  = "access_level"
      value = "restricted"
    },
    {
      name  = "backup_service"
      value = "false"
    }
  ]

  tags = [
    "service",
    "ssh",
    "management"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_service.test
  id = "123"

  identity = {
    custom_fields = [
      "service_owner:text",
      "monitoring_enabled:boolean",
      "access_level:text",
      "backup_service:boolean",
    ]
  }
}
