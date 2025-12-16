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
  role        = netbox_device_role.test.name
  site        = netbox_site.test.name
  status      = "active"
}

resource "netbox_interface" "test" {
  name   = "eth0"
  device = netbox_device.test.name
  type   = "1000base-t"
}

resource "netbox_tunnel" "test" {
  name          = "Test Tunnel"
  encapsulation = "ipsec-tunnel"
  status        = "active"
}

resource "netbox_tunnel_termination" "test" {
  tunnel             = netbox_tunnel.test.id
  role               = "peer"
  termination_type   = "dcim.interface"
  termination_id     = netbox_interface.test.id
  outside_ip_address = "1.2.3.4"
}
