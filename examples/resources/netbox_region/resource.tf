resource "netbox_region" "test" {
  name        = "Test Region"
  slug        = "test-region"
  description = "Geographic region for site grouping"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "region_code"
      value = "US-WEST"
    },
    {
      name  = "timezone"
      value = "America/Los_Angeles"
    },
    {
      name  = "contact_email"
      value = "us-west-ops@example.com"
    }
  ]

  tags = [
    "geographic-region",
    "us-west"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_region.test
  id = "123"

  identity = {
    custom_fields = [
      "region_code:text",
      "timezone:text",
      "contact_email:text",
    ]
  }
}
