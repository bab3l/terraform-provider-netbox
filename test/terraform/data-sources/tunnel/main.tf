# Tunnel Data Source Test

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

# Dependencies - create resources to test data sources
resource "netbox_tunnel" "test" {
  name          = "test-tunnel-ds"
  status        = "active"
  encapsulation = "gre"
  description   = "Test tunnel for data source"
}

# Test: Lookup tunnel by ID
data "netbox_tunnel" "by_id" {
  id = netbox_tunnel.test.id
}

# Test: Lookup tunnel by name
data "netbox_tunnel" "by_name" {
  name = netbox_tunnel.test.name

  depends_on = [netbox_tunnel.test]
}

# Output values for verification
output "by_id_name" {
  value = data.netbox_tunnel.by_id.name
}

output "by_id_status" {
  value = data.netbox_tunnel.by_id.status
}

output "by_id_encapsulation" {
  value = data.netbox_tunnel.by_id.encapsulation
}

output "by_name_id" {
  value = data.netbox_tunnel.by_name.id
}

output "by_id_description" {
  value = data.netbox_tunnel.by_id.description
}

# Validation output - all lookups should return the same ID
output "all_ids_match" {
  value = data.netbox_tunnel.by_id.id == data.netbox_tunnel.by_name.id
}
