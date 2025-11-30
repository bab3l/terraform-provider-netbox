# Outputs to verify tenant creation

output "basic_id" {
  description = "ID of the basic tenant"
  value       = netbox_tenant.basic.id
}

output "basic_name" {
  description = "Name of the basic tenant"
  value       = netbox_tenant.basic.name
}

output "basic_slug" {
  description = "Slug of the basic tenant"
  value       = netbox_tenant.basic.slug
}

output "complete_id" {
  description = "ID of the complete tenant"
  value       = netbox_tenant.complete.id
}

output "complete_name" {
  description = "Name of the complete tenant"
  value       = netbox_tenant.complete.name
}

output "complete_description" {
  description = "Description of the complete tenant"
  value       = netbox_tenant.complete.description
}

output "group_id" {
  description = "ID of the tenant group"
  value       = netbox_tenant_group.for_tenant.id
}

output "with_group_id" {
  description = "ID of the tenant with group"
  value       = netbox_tenant.with_group.id
}

output "with_group_group" {
  description = "Group ID of the tenant with group"
  value       = netbox_tenant.with_group.group
}

output "sibling1_id" {
  description = "ID of sibling tenant 1"
  value       = netbox_tenant.sibling1.id
}

output "sibling2_id" {
  description = "ID of sibling tenant 2"
  value       = netbox_tenant.sibling2.id
}

# Verification
output "group_assignment_valid" {
  description = "Verification that tenants are correctly assigned to group"
  value       = netbox_tenant.with_group.group == netbox_tenant_group.for_tenant.id && netbox_tenant.sibling1.group == netbox_tenant_group.for_tenant.id && netbox_tenant.sibling2.group == netbox_tenant_group.for_tenant.id
}
