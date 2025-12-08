# Outputs for ASN Range Data Source Test

# By ID outputs
output "by_id_id" {
  value       = data.netbox_asn_range.by_id.id
  description = "ID from ID lookup"
}

output "by_id_name" {
  value       = data.netbox_asn_range.by_id.name
  description = "Name from ID lookup"
}

output "by_id_slug" {
  value       = data.netbox_asn_range.by_id.slug
  description = "Slug from ID lookup"
}

output "by_id_rir" {
  value       = data.netbox_asn_range.by_id.rir
  description = "RIR from ID lookup"
}

output "by_id_start" {
  value       = data.netbox_asn_range.by_id.start
  description = "Start ASN from ID lookup"
}

output "by_id_end" {
  value       = data.netbox_asn_range.by_id.end
  description = "End ASN from ID lookup"
}

output "by_id_asn_count" {
  value       = data.netbox_asn_range.by_id.asn_count
  description = "ASN count from ID lookup"
}

# By name outputs
output "by_name_id" {
  value       = data.netbox_asn_range.by_name.id
  description = "ID from name lookup"
}

output "by_name_name" {
  value       = data.netbox_asn_range.by_name.name
  description = "Name from name lookup"
}

output "by_name_slug" {
  value       = data.netbox_asn_range.by_name.slug
  description = "Slug from name lookup"
}

# By slug outputs
output "by_slug_id" {
  value       = data.netbox_asn_range.by_slug.id
  description = "ID from slug lookup"
}

output "by_slug_name" {
  value       = data.netbox_asn_range.by_slug.name
  description = "Name from slug lookup"
}

output "by_slug_slug" {
  value       = data.netbox_asn_range.by_slug.slug
  description = "Slug from slug lookup"
}
