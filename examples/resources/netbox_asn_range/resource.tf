# Example: Creating ASN Ranges in Netbox
# ASN ranges allow you to define pools of AS Numbers for allocation

# First, create or reference a RIR (Regional Internet Registry)
resource "netbox_rir" "private" {
  name       = "Private ASN Range"
  slug       = "private-asn-range"
  is_private = true
}

# Basic ASN range with required fields only
resource "netbox_asn_range" "basic" {
  name  = "Private ASN Pool"
  slug  = "private-asn-pool"
  rir   = netbox_rir.private.id
  start = 64512 # Start of private ASN range
  end   = 65534 # End of private ASN range
}

# Full ASN range with all optional fields
resource "netbox_asn_range" "full" {
  name        = "Production ASN Pool"
  slug        = "production-asn-pool"
  rir         = netbox_rir.private.id
  start       = 64512
  end         = 64612
  tenant      = netbox_tenant.example.id
  description = "ASN range for production network devices"

  tags = {
    environment = "production"
    managed_by  = "terraform"
  }
}

# Reference an existing tenant
resource "netbox_tenant" "example" {
  name = "Example Tenant"
  slug = "example-tenant"
}
