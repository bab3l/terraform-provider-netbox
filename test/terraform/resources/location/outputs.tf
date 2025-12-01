# Outputs to verify location resource creation

output "basic_location_id" {
  description = "ID of the basic location"
  value       = netbox_location.basic.id
}

output "basic_location_name" {
  description = "Name of the basic location"
  value       = netbox_location.basic.name
}

output "basic_location_slug" {
  description = "Slug of the basic location"
  value       = netbox_location.basic.slug
}

output "basic_location_site" {
  description = "Site ID of the basic location"
  value       = netbox_location.basic.site
}

output "complete_location_id" {
  description = "ID of the complete location"
  value       = netbox_location.complete.id
}

output "complete_location_description" {
  description = "Description of the complete location"
  value       = netbox_location.complete.description
}

output "complete_location_status" {
  description = "Status of the complete location"
  value       = netbox_location.complete.status
}

output "complete_location_tenant" {
  description = "Tenant ID of the complete location"
  value       = netbox_location.complete.tenant
}

output "parent_location_id" {
  description = "ID of the parent location"
  value       = netbox_location.parent.id
}

output "child_location_id" {
  description = "ID of the child location"
  value       = netbox_location.child.id
}

output "child_location_parent" {
  description = "Parent ID of the child location"
  value       = netbox_location.child.parent
}

output "grandchild_location_id" {
  description = "ID of the grandchild location"
  value       = netbox_location.grandchild.id
}

output "grandchild_location_parent" {
  description = "Parent ID of the grandchild location"
  value       = netbox_location.grandchild.parent
}

output "planned_location_status" {
  description = "Status of the planned location"
  value       = netbox_location.planned.status
}

output "staging_location_status" {
  description = "Status of the staging location"
  value       = netbox_location.staging.status
}
