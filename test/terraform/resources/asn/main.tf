# ASN Resource Test

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

# First create the required RIR
resource "netbox_rir" "test" {
  name = "Test RIR for ASN"
  slug = "test-rir-asn"
}

# Test 1: Basic ASN creation
resource "netbox_asn" "basic" {
  asn = 65000
  rir = netbox_rir.test.id
}

# Test 2: ASN with all optional fields
resource "netbox_asn" "complete" {
  asn         = 65001
  rir         = netbox_rir.test.id
  description = "Test ASN for integration testing"
  comments    = "This ASN was created for integration testing."
}
