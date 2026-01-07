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
