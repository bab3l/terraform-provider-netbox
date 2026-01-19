resource "netbox_site" "example" {
  name        = "Example Site"
  slug        = "example-site"
  status      = "active"
  description = "An example site created with Terraform"
  facility    = "DC01"
  comments    = "This is a sample site configuration"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "datacenter_tier"
      value = "tier-3"
    },
    {
      name  = "power_capacity_kw"
      value = "500"
    },
    {
      name  = "cooling_type"
      value = "evaporative"
    },
    {
      name  = "primary_contact"
      value = "dc-ops@example.com"
    }
  ]

  tags = [
    "production",
    "datacenter"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_site.example
  id = "123"

  identity = {
    custom_fields = [
      "datacenter_tier:text",
      "power_capacity_kw:integer",
      "cooling_type:text",
      "primary_contact:text",
    ]
  }
}
