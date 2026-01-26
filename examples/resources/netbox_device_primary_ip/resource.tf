resource "netbox_site" "test" {
  name   = "Primary IP Site"
  slug   = "primary-ip-site"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Primary IP Manufacturer"
  slug = "primary-ip-manufacturer"
}

resource "netbox_device_type" "test" {
  model        = "Primary IP Device Type"
  slug         = "primary-ip-device-type"
  manufacturer = netbox_manufacturer.test.id
  u_height     = 1
}

resource "netbox_device_role" "test" {
  name  = "Primary IP Role"
  slug  = "primary-ip-role"
  color = "ff0000"
}

resource "netbox_device" "test" {
  name        = "primary-ip-device-1"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

resource "netbox_interface" "test" {
  device = netbox_device.test.id
  name   = "eth0"
  type   = "1000base-t"
}

resource "netbox_ip_address" "test" {
  address              = "10.0.0.10/24"
  status               = "active"
  assigned_object_type = "dcim.interface"
  assigned_object_id   = netbox_interface.test.id
}

resource "netbox_device_primary_ip" "test" {
  device      = netbox_device.test.id
  primary_ip4 = netbox_ip_address.test.id
}

import {
  to = netbox_device_primary_ip.test
  id = "123"
}
