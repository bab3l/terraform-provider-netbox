resource "netbox_fhrp_group" "test" {
  protocol  = "vrrp"
  group_id  = 10
  auth_type = "plaintext"
  auth_key  = "secret"
}

import {
  to = netbox_fhrp_group.test
  id = "123"
}
