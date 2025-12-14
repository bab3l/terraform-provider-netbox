# L2VPN Data Source Test

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

# Create an L2VPN resource to test the data source
resource "netbox_l2vpn" "test" {
  name        = "test-l2vpn-ds"
  slug        = "test-l2vpn-ds"
  type        = "vpls"
  identifier  = 20001
  description = "Test L2VPN for data source"
}

# Test 1: Look up by ID
data "netbox_l2vpn" "by_id" {
  id = netbox_l2vpn.test.id
}

# Test 2: Look up by slug
data "netbox_l2vpn" "by_slug" {
  slug = netbox_l2vpn.test.slug

  depends_on = [netbox_l2vpn.test]
}

# Test 3: Look up by name
data "netbox_l2vpn" "by_name" {
  name = netbox_l2vpn.test.name

  depends_on = [netbox_l2vpn.test]
}

# Outputs for verification
output "by_id_name" {
  value = data.netbox_l2vpn.by_id.name
}

output "by_id_slug" {
  value = data.netbox_l2vpn.by_id.slug
}

output "by_id_type" {
  value = data.netbox_l2vpn.by_id.type
}

output "by_id_identifier" {
  value = data.netbox_l2vpn.by_id.identifier
}

output "by_slug_name" {
  value = data.netbox_l2vpn.by_slug.name
}

output "by_name_slug" {
  value = data.netbox_l2vpn.by_name.slug
}

# Validation outputs
output "by_id_valid" {
  value = data.netbox_l2vpn.by_id.name == "test-l2vpn-ds"
}

output "by_slug_valid" {
  value = data.netbox_l2vpn.by_slug.name == "test-l2vpn-ds"
}

output "by_name_valid" {
  value = data.netbox_l2vpn.by_name.slug == "test-l2vpn-ds"
}

output "type_valid" {
  value = data.netbox_l2vpn.by_id.type == "vpls"
}

output "identifier_valid" {
  value = data.netbox_l2vpn.by_id.identifier == 20001
}
