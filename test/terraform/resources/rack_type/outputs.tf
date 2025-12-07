# Rack Type Outputs

# Basic rack type outputs
output "basic_id" {
  value = netbox_rack_type.basic.id
}

output "basic_model" {
  value = netbox_rack_type.basic.model
}

# Complete rack type outputs
output "complete_id" {
  value = netbox_rack_type.complete.id
}

output "complete_model" {
  value = netbox_rack_type.complete.model
}

output "complete_u_height" {
  value = netbox_rack_type.complete.u_height
}

output "complete_form_factor" {
  value = netbox_rack_type.complete.form_factor
}
