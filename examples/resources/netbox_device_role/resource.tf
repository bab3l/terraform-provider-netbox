resource "netbox_config_template" "test" {
  name          = "Test Config Template"
  template_code = "{{ device_role.name }}"
}

resource "netbox_device_role" "test" {
  name            = "Test Role"
  slug            = "test-role"
  color           = "ff0000"
  vm_role         = false
  description     = "Network switch role"
  config_template = netbox_config_template.test.id

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "role_category"
      value = "network"
    },
    {
      name  = "support_tier"
      value = "tier-1"
    },
    {
      name  = "monitoring_profile"
      value = "network-devices"
    }
  ]

  tags = [
    "network-role",
    "production"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_device_role.test
  id = "123"

  identity = {
    custom_fields = [
      "role_category:text",
      "support_tier:text",
      "monitoring_profile:text",
    ]
  }
}
