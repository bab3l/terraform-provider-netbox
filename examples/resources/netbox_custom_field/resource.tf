resource "netbox_custom_field" "test" {
  name          = "test_field"
  content_types = ["dcim.device"]
  type          = "text"
  label         = "Test Field"
  required      = false
}

import {
  to = netbox_custom_field.test
  id = "123"
}
