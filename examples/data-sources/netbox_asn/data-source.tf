# Look up ASN by ASN number
data "netbox_asn" "by_asn" {
  asn = "65001"
}

# Look up ASN by ID
data "netbox_asn" "by_id" {
  id = "123"
}

# Use ASN data in other resources
output "asn_number" {
  value = data.netbox_asn.by_id.asn
}

output "asn_rir" {
  value = data.netbox_asn.by_asn.rir
}

output "asn_tenant" {
  value = data.netbox_asn.by_id.tenant
}

output "asn_description" {
  value = data.netbox_asn.by_id.description
}

# Access all custom fields
output "asn_custom_fields" {
  value       = data.netbox_asn.by_id.custom_fields
  description = "All custom fields defined in NetBox for this ASN"
}

# Access specific custom fields by name
output "asn_organization" {
  value       = try([for cf in data.netbox_asn.by_id.custom_fields : cf.value if cf.name == "organization_name"][0], null)
  description = "Example: accessing a text custom field"
}

output "asn_public" {
  value       = try([for cf in data.netbox_asn.by_id.custom_fields : cf.value if cf.name == "is_public"][0], null)
  description = "Example: accessing a boolean custom field"
}
