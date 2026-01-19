resource "netbox_provider" "test" {
  name = "Test Provider"
  slug = "test-provider"
}

resource "netbox_provider_account" "test" {
  name             = "Test Account"
  account          = "1234567890"
  circuit_provider = netbox_provider.test.name
  description      = "Main provider account for datacenter services"
  comments         = "Primary account with billing contact"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "account_manager"
      value = "John Doe"
    },
    {
      name  = "billing_contact"
      value = "billing@example.com"
    },
    {
      name  = "annual_spend"
      value = "250000"
    }
  ]

  tags = [
    "production",
    "primary-account"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_provider_account.test
  id = "123"

  identity = {
    custom_fields = [
      "account_manager:text",
      "billing_contact:text",
      "annual_spend:integer",
    ]
  }
}
