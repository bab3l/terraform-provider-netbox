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
  manufacturer = netbox_manufacturer.test.id
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
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_virtual_device_context" "test" {
  name   = "Test VDC"
  device = netbox_device.test.id
  status = "active"
}

# Optional: seed owned custom fields during import
import {
  to = netbox_virtual_device_context.test
  id = "123"

  identity = {
    custom_fields = [
      "environment:text",
      "owner:text",
    ]
  }
}
