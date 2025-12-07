# RIR Outputs

# Basic RIR outputs
output "basic_id" {
  value = netbox_rir.basic.id
}

output "basic_name" {
  value = netbox_rir.basic.name
}

output "basic_slug" {
  value = netbox_rir.basic.slug
}

# Complete RIR outputs
output "complete_id" {
  value = netbox_rir.complete.id
}

output "complete_name" {
  value = netbox_rir.complete.name
}

output "complete_is_private" {
  value = netbox_rir.complete.is_private
}
