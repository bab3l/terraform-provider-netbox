# Outputs for contact assignment data source test

output "id_matches" {
  description = "Data source ID matches resource ID"
  value       = tostring(data.netbox_contact_assignment.by_id.id) == netbox_contact_assignment.test.id
}

output "object_type_matches" {
  description = "Data source object_type matches resource object_type"
  value       = data.netbox_contact_assignment.by_id.object_type == netbox_contact_assignment.test.object_type
}

output "object_id_matches" {
  description = "Data source object_id matches resource object_id"
  value       = data.netbox_contact_assignment.by_id.object_id == netbox_contact_assignment.test.object_id
}

output "contact_id_matches" {
  description = "Data source contact_id matches resource contact_id"
  value       = data.netbox_contact_assignment.by_id.contact_id == netbox_contact_assignment.test.contact_id
}

output "role_id_matches" {
  description = "Data source role_id matches resource role_id"
  value       = data.netbox_contact_assignment.by_id.role_id == netbox_contact_assignment.test.role_id
}

output "priority_matches" {
  description = "Data source priority matches resource priority"
  value       = data.netbox_contact_assignment.by_id.priority == netbox_contact_assignment.test.priority
}

output "contact_name" {
  description = "Contact name from data source"
  value       = data.netbox_contact_assignment.by_id.contact_name
}

output "role_name" {
  description = "Role name from data source"
  value       = data.netbox_contact_assignment.by_id.role_name
}
