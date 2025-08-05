# Tenant Group Resource Example

terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  server_url = "https://netbox.example.com"
  api_token  = "your-api-token-here"
}

# Create a top-level tenant group
resource "netbox_tenant_group" "corporate" {
  name        = "Corporate Tenants"
  slug        = "corporate"
  description = "Top-level group for all corporate tenants"

  tags = [
    {
      name = "corporate"
      slug = "corporate"
    }
  ]
}

# Create a child tenant group
resource "netbox_tenant_group" "subsidiaries" {
  name        = "Subsidiary Companies"
  slug        = "subsidiaries"
  parent      = netbox_tenant_group.corporate.id
  description = "Group for subsidiary company tenants"

  tags = [
    {
      name = "subsidiary"
      slug = "subsidiary"
    }
  ]
}

# Create another child tenant group
resource "netbox_tenant_group" "departments" {
  name        = "Departments"
  slug        = "departments"
  parent      = netbox_tenant_group.corporate.id
  description = "Group for internal department tenants"

  tags = [
    {
      name = "internal"
      slug = "internal"
    }
  ]
}

# Output the tenant group information
output "corporate_tenant_group" {
  value = {
    id          = netbox_tenant_group.corporate.id
    name        = netbox_tenant_group.corporate.name
    slug        = netbox_tenant_group.corporate.slug
    description = netbox_tenant_group.corporate.description
  }
}

output "subsidiaries_tenant_group" {
  value = {
    id          = netbox_tenant_group.subsidiaries.id
    name        = netbox_tenant_group.subsidiaries.name
    slug        = netbox_tenant_group.subsidiaries.slug
    parent      = netbox_tenant_group.subsidiaries.parent
    description = netbox_tenant_group.subsidiaries.description
  }
}
