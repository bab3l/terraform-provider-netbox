# Tag Data Source Test
# Tests the netbox_tag data source lookups by ID, name, and slug

terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  # Uses NETBOX_SERVER_URL and NETBOX_API_TOKEN environment variables
}

# Create a tag to look up
resource "netbox_tag" "test" {
  name        = "Tag DS Test"
  slug        = "tag-ds-test"
  color       = "9c27b0"
  description = "Tag for data source testing"
}

# Look up by ID
data "netbox_tag" "by_id" {
  id = netbox_tag.test.id
}

# Look up by name
data "netbox_tag" "by_name" {
  name = netbox_tag.test.name
}

# Look up by slug
data "netbox_tag" "by_slug" {
  slug = netbox_tag.test.slug
}

# Outputs for validation
output "id_lookup_matches" {
  value = data.netbox_tag.by_id.id == netbox_tag.test.id
}

output "name_lookup_matches" {
  value = data.netbox_tag.by_name.id == netbox_tag.test.id
}

output "slug_lookup_matches" {
  value = data.netbox_tag.by_slug.id == netbox_tag.test.id
}

output "by_id_name_valid" {
  value = data.netbox_tag.by_id.name == "Tag DS Test"
}

output "by_id_slug_valid" {
  value = data.netbox_tag.by_id.slug == "tag-ds-test"
}

output "by_id_color_valid" {
  value = data.netbox_tag.by_id.color == "9c27b0"
}

output "by_id_description_valid" {
  value = data.netbox_tag.by_id.description == "Tag for data source testing"
}
