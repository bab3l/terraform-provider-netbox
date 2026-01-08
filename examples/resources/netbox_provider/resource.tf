resource "netbox_provider" "test" {
  name     = "Test Provider"
  slug     = "test-provider"
  comments = "Primary internet service provider"

  # Partial custom fields management
  # Only the custom fields specified here are managed by Terraform
  # Other custom fields set in NetBox (e.g., via UI) are preserved
  custom_fields = [
    {
      name  = "account_manager"
      value = "Jane Smith"
    },
    {
      name  = "support_phone"
      value = "+1-800-555-0123"
    },
    {
      name  = "support_email"
      value = "support@testprovider.com"
    }
  ]

  tags = [
    "tier1-provider",
    "production"
  ]
}
