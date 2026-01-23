resource "netbox_export_template" "test" {
  name           = "Test Export Template"
  template_code  = "{% for device in queryset %}{{ device.name }}{% endfor %}"
  content_type   = "dcim.device"
  file_extension = "txt"
  mime_type      = "text/plain"
}

import {
  to = netbox_export_template.test
  id = "123"
}
