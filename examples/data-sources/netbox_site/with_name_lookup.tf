# Site Data Source with Name Lookup Example

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

# Get information about a site by ID
data "netbox_site" "example_by_id" {
  id = "1"
}

# Get information about a site by slug
data "netbox_site" "example_by_slug" {
  slug = "dc-east-1"
}

# Get information about a site by name (NEW!)
data "netbox_site" "example_by_name" {
  name = "Primary Datacenter East"
}

# Output all three lookup methods
output "site_by_id" {
  value = {
    id          = data.netbox_site.example_by_id.id
    name        = data.netbox_site.example_by_id.name
    slug        = data.netbox_site.example_by_id.slug
    status      = data.netbox_site.example_by_id.status
    description = data.netbox_site.example_by_id.description
  }
}

output "site_by_slug" {
  value = {
    id          = data.netbox_site.example_by_slug.id
    name        = data.netbox_site.example_by_slug.name
    slug        = data.netbox_site.example_by_slug.slug
    status      = data.netbox_site.example_by_slug.status
    description = data.netbox_site.example_by_slug.description
  }
}

output "site_by_name" {
  value = {
    id          = data.netbox_site.example_by_name.id
    name        = data.netbox_site.example_by_name.name
    slug        = data.netbox_site.example_by_name.slug
    status      = data.netbox_site.example_by_name.status
    description = data.netbox_site.example_by_name.description
  }
}

# Example showing lookup preference order
# ID is most efficient (direct lookup)
# Slug is second best (unique and URL-friendly)
# Name is useful for human-readable configs but may not be unique

locals {
  lookup_methods = {
    by_id = {
      method     = "ID"
      efficiency = "highest"
      uniqueness = "guaranteed"
      example    = "data.netbox_site.example { id = \"1\" }"
    }
    by_slug = {
      method     = "Slug"
      efficiency = "high"
      uniqueness = "guaranteed"
      example    = "data.netbox_site.example { slug = \"dc-east-1\" }"
    }
    by_name = {
      method     = "Name"
      efficiency = "medium"
      uniqueness = "not_guaranteed"
      example    = "data.netbox_site.example { name = \"Primary Datacenter\" }"
      note       = "May return error if multiple sites have the same name"
    }
  }
}

output "lookup_comparison" {
  description = "Comparison of different lookup methods"
  value       = local.lookup_methods
}
