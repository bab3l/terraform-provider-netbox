# ASN Data Source Test

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
resource "netbox_rir" "test" {
  name = "Test RIR for ASN DS"
  slug = "test-rir-asn-ds"
}

resource "netbox_asn" "test" {
  asn         = 65001
  rir         = netbox_rir.test.id
  description = "Test ASN for data source"
}

# Test: Lookup ASN by ID
data "netbox_asn" "by_id" {
  id = netbox_asn.test.id
}

# Test: Lookup ASN by asn number
data "netbox_asn" "by_asn" {
  asn = netbox_asn.test.asn

  depends_on = [netbox_asn.test]
}
