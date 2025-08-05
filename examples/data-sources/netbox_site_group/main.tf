# Site Group Data Source Example

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

# Get information about a site group by ID
data "netbox_site_group" "example_by_id" {
  id = "1"
}

# Get information about a site group by slug
data "netbox_site_group" "example_by_slug" {
  slug = "north-america"
}

# Get information about a site group by name
data "netbox_site_group" "example_by_name" {
  name = "North America"
}

# Output the site group information
output "site_group_by_id" {
  value = {
    id          = data.netbox_site_group.example_by_id.id
    name        = data.netbox_site_group.example_by_id.name
    slug        = data.netbox_site_group.example_by_id.slug
    parent      = data.netbox_site_group.example_by_id.parent
    description = data.netbox_site_group.example_by_id.description
  }
}

output "site_group_by_slug" {
  value = {
    id          = data.netbox_site_group.example_by_slug.id
    name        = data.netbox_site_group.example_by_slug.name
    slug        = data.netbox_site_group.example_by_slug.slug
    parent      = data.netbox_site_group.example_by_slug.parent
    description = data.netbox_site_group.example_by_slug.description
  }
}

output "site_group_by_name" {
  value = {
    id          = data.netbox_site_group.example_by_name.id
    name        = data.netbox_site_group.example_by_name.name
    slug        = data.netbox_site_group.example_by_name.slug
    parent      = data.netbox_site_group.example_by_name.parent
    description = data.netbox_site_group.example_by_name.description
  }
}
