# Tunnel Group Data Source Test

terraform {
  required_version = ">= 1.0"
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  # Uses NETBOX_SERVER_URL and NETBOX_API_TOKEN environment variables
}

# Dependencies - create resources to test data sources
resource "netbox_tunnel_group" "test" {
  name        = "test-tunnel-group-ds"
  slug        = "test-tunnel-group-ds"
  description = "Test tunnel group for data source"
}

# Test: Lookup tunnel group by ID
data "netbox_tunnel_group" "by_id" {
  id = netbox_tunnel_group.test.id
}

# Test: Lookup tunnel group by name
data "netbox_tunnel_group" "by_name" {
  name = netbox_tunnel_group.test.name

  depends_on = [netbox_tunnel_group.test]
}

# Test: Lookup tunnel group by slug
data "netbox_tunnel_group" "by_slug" {
  slug = netbox_tunnel_group.test.slug

  depends_on = [netbox_tunnel_group.test]
}

# Output values for verification
output "by_id_name" {
  value = data.netbox_tunnel_group.by_id.name
}

output "by_id_slug" {
  value = data.netbox_tunnel_group.by_id.slug
}

output "by_name_id" {
  value = data.netbox_tunnel_group.by_name.id
}

output "by_slug_id" {
  value = data.netbox_tunnel_group.by_slug.id
}

output "by_id_description" {
  value = data.netbox_tunnel_group.by_id.description
}

# Validation output - all lookups should return the same ID
output "all_ids_match" {
  value = (data.netbox_tunnel_group.by_id.id == data.netbox_tunnel_group.by_name.id &&
           data.netbox_tunnel_group.by_name.id == data.netbox_tunnel_group.by_slug.id)
}
