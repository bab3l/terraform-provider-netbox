# Outputs for contact assignment resource test

output "basic_id" {
  description = "ID of the basic contact assignment"
  value       = netbox_contact_assignment.basic.id
}

output "basic_id_valid" {
  description = "Basic contact assignment has valid ID"
  value       = netbox_contact_assignment.basic.id != ""
}

output "basic_object_type" {
  description = "Object type of the basic contact assignment"
  value       = netbox_contact_assignment.basic.object_type
}

output "basic_contact_id" {
  description = "Contact ID of the basic contact assignment"
  value       = netbox_contact_assignment.basic.contact_id
}

output "basic_role_id" {
  description = "Role ID of the basic contact assignment"
  value       = netbox_contact_assignment.basic.role_id
}

output "with_role_id" {
  description = "ID of the contact assignment with role"
  value       = netbox_contact_assignment.with_role.id
}

output "with_role_priority" {
  description = "Priority of the contact assignment with role"
  value       = netbox_contact_assignment.with_role.priority
}

output "with_role_role_id" {
  description = "Role ID of the contact assignment with role"
  value       = netbox_contact_assignment.with_role.role_id
}

output "site_id" {
  description = "ID of the test site"
  value       = netbox_site.test.id
}

output "contact_id" {
  description = "ID of the test contact"
  value       = netbox_contact.test.id
}

output "role_id" {
  description = "ID of the test contact role"
  value       = netbox_contact_role.test.id
}
