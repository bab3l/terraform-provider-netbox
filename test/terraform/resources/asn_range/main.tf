# ASN Range Resource Test
# Tests creation and management of ASN ranges in Netbox

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
  name        = "Test RIR for ASN Range TF"
  slug        = "test-rir-asn-range-tf"
  is_private  = false
  description = "RIR for ASN range testing"
}

# Create tenant for ASN range (optional)
resource "netbox_tenant" "test" {
  name        = "Test Tenant for ASN Range TF"
  slug        = "test-tenant-asn-range-tf"
  description = "Tenant for ASN range testing"
}

# Create basic ASN range
resource "netbox_asn_range" "basic" {
  name  = "Basic ASN Range"
  slug  = "basic-asn-range"
  rir   = netbox_rir.test.id
  start = 64512
  end   = 64612
}

# Create ASN range with all optional fields
resource "netbox_asn_range" "full" {
  name        = "Full ASN Range"
  slug        = "full-asn-range"
  rir         = netbox_rir.test.id
  start       = 65000
  end         = 65100
  tenant      = netbox_tenant.test.id
  description = "Full ASN range with all optional fields"
}
