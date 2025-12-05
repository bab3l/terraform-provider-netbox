# Tag Resource Test
# Tests the netbox_tag resource CRUD operations

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

# Basic tag with required fields only
resource "netbox_tag" "basic" {
  name = "Tag Resource Test Basic"
  slug = "tag-resource-test-basic"
}

# Tag with all optional fields
resource "netbox_tag" "full" {
  name        = "Tag Resource Test Full"
  slug        = "tag-resource-test-full"
  color       = "ff5722"
  description = "A test tag with all fields populated"
}

# Tag with object type restrictions
resource "netbox_tag" "restricted" {
  name         = "Tag Resource Test Restricted"
  slug         = "tag-resource-test-restricted"
  color        = "2196f3"
  description  = "A tag restricted to specific object types"
  object_types = ["dcim.device", "dcim.site"]
}

# Outputs for validation
output "basic_id_valid" {
  value = can(tonumber(netbox_tag.basic.id))
}

output "basic_name_valid" {
  value = netbox_tag.basic.name == "Tag Resource Test Basic"
}

output "basic_slug_valid" {
  value = netbox_tag.basic.slug == "tag-resource-test-basic"
}

output "full_color_valid" {
  value = netbox_tag.full.color == "ff5722"
}

output "full_description_valid" {
  value = netbox_tag.full.description == "A test tag with all fields populated"
}

output "restricted_object_types_valid" {
  value = length(netbox_tag.restricted.object_types) == 2
}
