resource "netbox_rir" "test" {
  name        = "Test RIR"
  slug        = "test-rir"
  description = "Regional Internet Registry for testing"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "registry_region"
      value = "North America"
    },
    {
      name  = "contact_email"
      value = "admin@test-rir.org"
    },
    {
      name  = "whois_server"
      value = "whois.test-rir.org"
    }
  ]

  tags = [
    "rir",
    "test-registry"
  ]
}
