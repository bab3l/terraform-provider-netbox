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
