resource "netbox_contact" "test" {
  name = "John Doe"
}

resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
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

resource "netbox_device_role" "test" {
  name    = "Test Role"
  slug    = "test-role"
  color   = "ff0000"
  vm_role = false
}

resource "netbox_device" "test" {
  name        = "Test Device"
  device_type = netbox_device_type.test.slug
  role        = netbox_device_role.test.slug
  site        = netbox_site.test.slug
}

resource "netbox_contact_role" "test" {
  name = "Admin"
  slug = "admin"
}

resource "netbox_contact_assignment" "test" {
  content_type = "dcim.device"
  object_id    = netbox_device.test.id
  contact      = netbox_contact.test.name
  role         = netbox_contact_role.test.slug
  priority     = "primary"
  description  = "Primary contact for device support"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "notification_method"
      value = "email"
    },
    {
      name  = "availability_hours"
      value = "24x7"
    },
    {
      name  = "escalation_level"
      value = "1"
    }
  ]

  tags = [
    "contact-assignment",
    "primary"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_contact_assignment.test
  id = "123"

  identity = {
    custom_fields = [
      "notification_method:text",
      "availability_hours:text",
      "escalation_level:integer",
    ]
  }
}
