# ASN Range Data Source Test
# Tests retrieving ASN ranges from Netbox by various lookup methods

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

# Create RIR for ASN range
resource "netbox_rir" "test" {
  name        = "Test RIR for ASN Range DS"
  slug        = "test-rir-asn-range-ds"
  is_private  = false
  description = "RIR for ASN range data source testing"
}

# Create ASN range to lookup
resource "netbox_asn_range" "test" {
  name        = "Test ASN Range for DS"
  slug        = "test-asn-range-for-ds"
  rir         = netbox_rir.test.id
  start       = 64700
  end         = 64800
  description = "ASN range for data source testing"
}

# Look up by ID
data "netbox_asn_range" "by_id" {
  id = netbox_asn_range.test.id
}

# Look up by name
data "netbox_asn_range" "by_name" {
  name = netbox_asn_range.test.name
}

# Look up by slug
data "netbox_asn_range" "by_slug" {
  slug = netbox_asn_range.test.slug
}
