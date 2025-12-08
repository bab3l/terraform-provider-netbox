# Example: Looking up ASN Ranges from Netbox

# Look up an ASN range by ID
data "netbox_asn_range" "by_id" {
  id = 1
}

# Look up an ASN range by name
data "netbox_asn_range" "by_name" {
  name = "Private ASN Pool"
}

# Look up an ASN range by slug
data "netbox_asn_range" "by_slug" {
  slug = "private-asn-pool"
}

# Use the data source to reference ASN range properties
output "asn_range_details" {
  value = {
    id        = data.netbox_asn_range.by_name.id
    name      = data.netbox_asn_range.by_name.name
    slug      = data.netbox_asn_range.by_name.slug
    rir       = data.netbox_asn_range.by_name.rir
    start     = data.netbox_asn_range.by_name.start
    end       = data.netbox_asn_range.by_name.end
    asn_count = data.netbox_asn_range.by_name.asn_count
  }
}
