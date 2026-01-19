resource "netbox_role" "test" {
  name        = "Test Role"
  slug        = "test-role"
  description = "Example IPAM role for network infrastructure"
  weight      = 1000

  # Partial custom fields management
  # Only specified custom fields are managed by Terraform
  # Other custom fields set in NetBox UI or other tools are preserved
  custom_fields = [
    {
      name  = "cost_center"
      type  = "text"
      value = "IT-NETWORK-001"
    },
    {
      name  = "managed_by_terraform"
      type  = "boolean"
      value = "true"
    }
  ]

  tags = [
    {
      name = "production"
      slug = "production"
    }
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_role.test
  id = "123"

  identity = {
    custom_fields = [
      "cost_center:text",
      "managed_by_terraform:boolean",
    ]
  }
}
