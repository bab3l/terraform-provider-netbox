resource "netbox_inventory_item_role" "test" {
  name  = "Test Role"
  slug  = "test-role"
  color = "ff0000"
}

import {
  to = netbox_inventory_item_role.test
  id = "123"
}
