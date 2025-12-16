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
  model           = "Test Model"
  slug            = "test-model"
  manufacturer_id = netbox_manufacturer.test.id
  u_height        = 1
}

resource "netbox_device_role" "test" {
  name    = "Test Role"
  slug    = "test-role"
  color   = "ff0000"
  vm_role = false
}

resource "netbox_device" "test" {
  name           = "Test Device"
  device_type_id = netbox_device_type.test.id
  role_id        = netbox_device_role.test.id
  site_id        = netbox_site.test.id
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
}
