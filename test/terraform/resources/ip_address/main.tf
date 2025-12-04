# IP Address Integration Test
# Tests the netbox_ip_address resource with basic and complete configurations

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

# Prerequisites
resource "netbox_vrf" "test" {
  name        = "IP Address Test VRF"
  rd          = "65000:400"
  description = "VRF for IP address testing"
}

resource "netbox_tenant" "test" {
  name = "IP Address Test Tenant"
  slug = "ip-address-test-tenant"
}

# Basic IP Address with only required fields
resource "netbox_ip_address" "basic" {
  address = "10.100.0.1/24"
}

# Complete IP Address with all optional fields
resource "netbox_ip_address" "complete" {
  address     = "10.100.1.1/24"
  status      = "active"
  description = "Complete IP address for integration testing"
  comments    = "Created by terraform integration test"
  vrf         = netbox_vrf.test.id
  tenant      = netbox_tenant.test.id
  dns_name    = "test-server.example.com"
}

# IPv6 Address
resource "netbox_ip_address" "ipv6" {
  address     = "2001:db8::1/64"
  status      = "active"
  description = "IPv6 address test"
}

# Reserved IP Address
resource "netbox_ip_address" "reserved" {
  address     = "10.100.2.1/24"
  status      = "reserved"
  description = "Reserved IP address test"
}

# DHCP IP Address
resource "netbox_ip_address" "dhcp" {
  address     = "10.100.3.1/24"
  status      = "dhcp"
  description = "DHCP IP address test"
}

# SLAAC IP Address
resource "netbox_ip_address" "slaac" {
  address     = "2001:db8::100/64"
  status      = "slaac"
  description = "SLAAC IPv6 address test"
}
