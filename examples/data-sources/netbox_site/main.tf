# Site Data Source Example

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

# Output the site information
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
