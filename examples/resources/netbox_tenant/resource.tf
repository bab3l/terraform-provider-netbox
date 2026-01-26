terraform {
  required_providers {
    netbox = {
      source = "local/bab3l/netbox"
    }
  }
}

provider "netbox" {
  server_url = "https://netbox.example.com"
  api_token  = "your_api_token_here"
}

# Create a tenant group first
resource "netbox_tenant_group" "example_group" {
  name        = "Example Tenant Group"
  slug        = "example-tenant-group"
  description = "A group for organizing tenants"
}

# Create a tenant within the group
resource "netbox_tenant" "example_tenant" {
  name        = "Example Tenant"
  slug        = "example-tenant"
  group       = netbox_tenant_group.example_group.id
  description = "An example tenant"
  comments    = "This tenant is used for demonstration purposes"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "account_id"
      value = "ACCT-12345"
    },
    {
      name  = "billing_contact"
      value = "billing@example.com"
    },
    {
      name  = "contract_end_date"
      value = "2026-12-31"
    },
    {
      name  = "support_tier"
      value = "premium"
    }
  ]

  tags = [
    {
      name = "production"
      slug = "production"
    },
    {
      name = "critical"
      slug = "critical"
    }
  ]
}

# Create a tenant without a group
resource "netbox_tenant" "standalone_tenant" {
  name        = "Standalone Tenant"
  slug        = "standalone-tenant"
  description = "A tenant without a group"
}

# Optional: seed owned custom fields during import
import {
  to = netbox_tenant.example_tenant
  id = "123"

  identity = {
    custom_fields = [
      "account_id:text",
      "billing_contact:text",
      "contract_end_date:date",
      "support_tier:text",
    ]
  }
}
