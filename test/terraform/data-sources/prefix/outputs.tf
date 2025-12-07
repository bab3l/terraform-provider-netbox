# Prefix Data Source Outputs

output "by_id_prefix" {
  value = data.netbox_prefix.by_id.prefix
}

output "by_id_status" {
  value = data.netbox_prefix.by_id.status
}

output "by_prefix_id" {
  value = data.netbox_prefix.by_prefix.id
}

output "by_prefix_description" {
  value = data.netbox_prefix.by_prefix.description
}
