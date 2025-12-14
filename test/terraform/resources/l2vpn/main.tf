# L2VPN Resource Test

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

# Test 1: Basic L2VPN with required fields only
resource "netbox_l2vpn" "basic" {
  name = "test-l2vpn-basic"
  slug = "test-l2vpn-basic"
  type = "vpls"
}

# Test 2: L2VPN with all optional fields
resource "netbox_l2vpn" "complete" {
  name        = "test-l2vpn-complete"
  slug        = "test-l2vpn-complete"
  type        = "vxlan"
  identifier  = 10001
  description = "Complete L2VPN for integration testing"
  comments    = "This L2VPN was created for integration testing."
}

# Test 3: L2VPN with different types
resource "netbox_l2vpn" "evpn" {
  name = "test-l2vpn-evpn"
  slug = "test-l2vpn-evpn"
  type = "vxlan-evpn"
}

resource "netbox_l2vpn" "vpws" {
  name = "test-l2vpn-vpws"
  slug = "test-l2vpn-vpws"
  type = "vpws"
}

# Test 4: Output values for verification
output "basic_l2vpn_id" {
  value = netbox_l2vpn.basic.id
}

output "basic_l2vpn_name" {
  value = netbox_l2vpn.basic.name
}

output "basic_l2vpn_slug" {
  value = netbox_l2vpn.basic.slug
}

output "basic_l2vpn_type" {
  value = netbox_l2vpn.basic.type
}

output "complete_l2vpn_name" {
  value = netbox_l2vpn.complete.name
}

output "complete_l2vpn_description" {
  value = netbox_l2vpn.complete.description
}

output "complete_l2vpn_identifier" {
  value = netbox_l2vpn.complete.identifier
}

# Validation outputs
output "basic_id_valid" {
  value = netbox_l2vpn.basic.id != "" && netbox_l2vpn.basic.id != null
}

output "basic_name_valid" {
  value = netbox_l2vpn.basic.name == "test-l2vpn-basic"
}

output "basic_slug_valid" {
  value = netbox_l2vpn.basic.slug == "test-l2vpn-basic"
}

output "basic_type_valid" {
  value = netbox_l2vpn.basic.type == "vpls"
}

output "complete_identifier_valid" {
  value = netbox_l2vpn.complete.identifier == 10001
}
