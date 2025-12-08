# Outputs for ASN Range Resource Test

output "basic_asn_range_id" {
  value       = netbox_asn_range.basic.id
  description = "ID of the basic ASN range"
}

output "basic_asn_range_name" {
  value       = netbox_asn_range.basic.name
  description = "Name of the basic ASN range"
}

output "basic_asn_range_slug" {
  value       = netbox_asn_range.basic.slug
  description = "Slug of the basic ASN range"
}

output "basic_asn_range_start" {
  value       = netbox_asn_range.basic.start
  description = "Start ASN of the basic range"
}

output "basic_asn_range_end" {
  value       = netbox_asn_range.basic.end
  description = "End ASN of the basic range"
}

output "full_asn_range_id" {
  value       = netbox_asn_range.full.id
  description = "ID of the full ASN range"
}

output "full_asn_range_description" {
  value       = netbox_asn_range.full.description
  description = "Description of the full ASN range"
}

output "full_asn_range_tenant" {
  value       = netbox_asn_range.full.tenant
  description = "Tenant ID of the full ASN range"
}

output "full_asn_range_tags" {
  value       = netbox_asn_range.full.tags
  description = "Tags of the full ASN range"
}
