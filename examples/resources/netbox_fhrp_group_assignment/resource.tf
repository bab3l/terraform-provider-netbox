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

resource "netbox_interface" "test" {
  name   = "GigabitEthernet1/0/1"
  device = netbox_device.test.id
  type   = "1000base-t"
}

resource "netbox_fhrp_group" "test" {
  protocol = "vrrp"
  group_id = 10
}

resource "netbox_fhrp_group_assignment" "test" {
  group_id       = netbox_fhrp_group.test.id
  interface_type = "dcim.interface"
  interface_id   = netbox_interface.test.id
  priority       = 100
}

import {
  to = netbox_fhrp_group_assignment.test
  id = "123"
}
