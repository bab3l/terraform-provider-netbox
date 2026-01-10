# Example: Looking up ASN Ranges from Netbox

# Look up an ASN range by ID
data "netbox_asn_range" "by_id" {
  id = "1"
}

# Look up an ASN range by name
data "netbox_asn_range" "by_name" {
  name = "Private ASN Pool"
}

# Look up an ASN range by slug
data "netbox_asn_range" "by_slug" {
  slug = "private-asn-pool"
}

# Use ASN range data in other resources
output "asn_range_name" {
  value = data.netbox_asn_range.by_name.name
}

output "asn_range_slug" {
  value = data.netbox_asn_range.by_slug.slug
}

output "asn_range_rir" {
  value = data.netbox_asn_range.by_id.rir
}

output "asn_range_start" {
  value = data.netbox_asn_range.by_id.start
}

output "asn_range_end" {
  value = data.netbox_asn_range.by_id.end
}

output "asn_range_count" {
  value = data.netbox_asn_range.by_name.asn_count
}

# Access all custom fields
output "asn_range_custom_fields" {
  value       = data.netbox_asn_range.by_id.custom_fields
  description = "All custom fields defined in NetBox for this ASN range"
}

# Access specific custom fields by name
output "asn_range_purpose" {
  value       = try([for cf in data.netbox_asn_range.by_id.custom_fields : cf.value if cf.name == "purpose"][0], null)
  description = "Example: accessing a select custom field"
}

output "asn_range_contact" {
  value       = try([for cf in data.netbox_asn_range.by_id.custom_fields : cf.value if cf.name == "contact_email"][0], null)
  description = "Example: accessing a text custom field"
}
