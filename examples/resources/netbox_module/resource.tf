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

resource "netbox_module_bay" "test" {
  name   = "Module Bay 1"
  device = netbox_device.test.name
}

resource "netbox_module_type" "test" {
  model        = "Test Module Type"
  manufacturer = netbox_manufacturer.test.name
}

resource "netbox_module" "test" {
  device      = netbox_device.test.name
  module_bay  = netbox_module_bay.test.name
  module_type = netbox_module_type.test.model
  status      = "active"
}
