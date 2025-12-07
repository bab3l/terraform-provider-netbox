# Role Outputs

# Basic role outputs
output "basic_id" {
  value = netbox_role.basic.id
}

output "basic_name" {
  value = netbox_role.basic.name
}

output "basic_slug" {
  value = netbox_role.basic.slug
}

# Complete role outputs
output "complete_id" {
  value = netbox_role.complete.id
}

output "complete_name" {
  value = netbox_role.complete.name
}

output "complete_weight" {
  value = netbox_role.complete.weight
}
