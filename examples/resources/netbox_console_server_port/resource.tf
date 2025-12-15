resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_device_role" "test" {
  name  = "Console Server Role"
  slug  = "console-server-role"
  color = "0000ff"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_type" "test" {
  model        = "Console Server Model"
  slug         = "console-server-model"
  manufacturer = netbox_manufacturer.test.id
  u_height     = 1
}

resource "netbox_device" "test" {
  name        = "test-console-server-1"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

resource "netbox_console_server_port" "test" {
  name   = "Port 1"
  device = netbox_device.test.id
  type   = "rj-45"
}
