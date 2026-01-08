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
  manufacturer = netbox_manufacturer.test.slug
  u_height     = 1
}

resource "netbox_device" "test" {
  name        = "Test Device"
  device_type = netbox_device_type.test.model
  role        = netbox_device_role.test.slug
  site        = netbox_site.test.slug
}

resource "netbox_service" "test" {
  name        = "SSH"
  protocol    = "tcp"
  ports       = [22]
  device      = netbox_device.test.name
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
