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

# Look up a tenant by ID
data "netbox_tenant" "by_id" {
  id = "123"
}

# Look up a tenant by slug
data "netbox_tenant" "by_slug" {
  slug = "example-tenant"
}

# Look up a tenant by name
data "netbox_tenant" "by_name" {
  name = "Example Tenant"
}

# Output the tenant information
output "tenant_info" {
  value = {
    id          = data.netbox_tenant.by_slug.id
    name        = data.netbox_tenant.by_slug.name
    slug        = data.netbox_tenant.by_slug.slug
    group       = data.netbox_tenant.by_slug.group
    description = data.netbox_tenant.by_slug.description
    comments    = data.netbox_tenant.by_slug.comments
  }
}
