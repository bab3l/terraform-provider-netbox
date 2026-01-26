resource "netbox_manufacturer" "test" {
  name = "Cisco"
  slug = "cisco"
}

resource "netbox_platform" "test" {
  name         = "Cisco IOS"
  slug         = "cisco-ios"
  manufacturer = netbox_manufacturer.test.id
}

import {
  to = netbox_platform.test
  id = "123"
}
