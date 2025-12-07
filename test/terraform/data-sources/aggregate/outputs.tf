# Aggregate Data Source Outputs

output "by_id_prefix" {
  value = data.netbox_aggregate.by_id.prefix
}

output "by_id_rir" {
  value = data.netbox_aggregate.by_id.rir
}

output "by_prefix_id" {
  value = data.netbox_aggregate.by_prefix.id
}

output "by_prefix_description" {
  value = data.netbox_aggregate.by_prefix.description
}
