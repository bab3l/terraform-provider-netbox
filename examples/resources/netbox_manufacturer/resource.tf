resource "netbox_manufacturer" "test" {
  name        = "Test Manufacturer"
  slug        = "test-manufacturer"
  description = "Network equipment manufacturer"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "support_level"
      value = "premium"
    },
    {
      name  = "support_phone"
      value = "+1-800-555-0100"
    },
    {
      name  = "warranty_default_years"
      value = "3"
    }
  ]

  tags = [
    "vendor",
    "network-equipment"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_manufacturer.test
  id = "123"

  identity = {
    custom_fields = [
      "support_level:text",
      "support_phone:text",
      "warranty_default_years:integer",
    ]
  }
}
