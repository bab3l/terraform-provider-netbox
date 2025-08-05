# Tenant Group Data Source with Hierarchy Example

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

# Create a parent tenant group
resource "netbox_tenant_group" "enterprise" {
  name        = "Enterprise Clients"
  slug        = "enterprise"
  description = "Top-level group for enterprise clients"
}

# Create a child tenant group
resource "netbox_tenant_group" "fortune_500" {
  name        = "Fortune 500"
  slug        = "fortune-500"
  parent      = netbox_tenant_group.enterprise.id
  description = "Fortune 500 company tenants"
}

# Use data source to read information about the parent group
data "netbox_tenant_group" "parent_info" {
  id = netbox_tenant_group.enterprise.id
}

# Use data source to read information about the child group
data "netbox_tenant_group" "child_info" {
  slug = netbox_tenant_group.fortune_500.slug
}

# Output showing the hierarchical relationship
output "hierarchy_example" {
  value = {
    parent = {
      id   = data.netbox_tenant_group.parent_info.id
      name = data.netbox_tenant_group.parent_info.name
      slug = data.netbox_tenant_group.parent_info.slug
    }
    child = {
      id     = data.netbox_tenant_group.child_info.id
      name   = data.netbox_tenant_group.child_info.name
      slug   = data.netbox_tenant_group.child_info.slug
      parent = data.netbox_tenant_group.child_info.parent
    }
  }
}
