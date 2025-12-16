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
  group       = netbox_tenant_group.example_group.slug
  description = "An example tenant"
  comments    = "This tenant is used for demonstration purposes"

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

  custom_fields = [
    {
      name  = "cost_center"
      type  = "text"
      value = "CC-1234"
    },
    {
      name  = "contact_email"
      type  = "text"
      value = "admin@example.com"
    }
  ]
}

# Create a tenant without a group
resource "netbox_tenant" "standalone_tenant" {
  name        = "Standalone Tenant"
  slug        = "standalone-tenant"
  description = "A tenant without a group"
}
