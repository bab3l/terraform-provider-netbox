resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_device_role" "test" {
  name  = "Patch Panel Role"
  slug  = "patch-panel-role"
  color = "cccccc"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_type" "test" {
  model        = "Patch Panel Model"
  slug         = "patch-panel-model"
  manufacturer = netbox_manufacturer.test.id
  u_height     = 1
}

resource "netbox_device" "test" {
  name        = "test-patch-panel-1"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

resource "netbox_rear_port" "test" {
  name      = "Rear Port 1"
  device    = netbox_device.test.id
  type      = "8p8c"
  positions = 1
}
