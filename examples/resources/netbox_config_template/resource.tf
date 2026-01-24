resource "netbox_config_template" "test" {
  name          = "Test Config Template"
  template_code = "hostname {{ device.name }}"
  environment_params = jsonencode({
    "foo" : "bar"
  })
}

import {
  to = netbox_config_template.test
  id = "123"
}
