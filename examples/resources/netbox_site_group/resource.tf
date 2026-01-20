resource "netbox_site_group" "example" {
  name        = "Example Site Group"
  slug        = "example-site-group"
  description = "An example site group created with Terraform"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "business_region"
      value = "North America"
    },
    {
      name  = "regional_manager"
      value = "Jane Doe"
    },
    {
      name  = "cost_center"
      value = "DC-NA-001"
    }
  ]

  tags = [
    "site-group",
    "north-america"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_site_group.example
  id = "123"

  identity = {
    custom_fields = [
      "business_region:text",
      "regional_manager:text",
      "cost_center:text",
    ]
  }
}
