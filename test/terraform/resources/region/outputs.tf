# Outputs to verify resource creation

output "basic_region_id" {
  description = "ID of the basic region"
  value       = netbox_region.basic.id
}

output "basic_region_name" {
  description = "Name of the basic region"
  value       = netbox_region.basic.name
}

output "basic_region_slug" {
  description = "Slug of the basic region"
  value       = netbox_region.basic.slug
}

output "complete_region_id" {
  description = "ID of the complete region"
  value       = netbox_region.complete.id
}

output "complete_region_description" {
  description = "Description of the complete region"
  value       = netbox_region.complete.description
}

output "parent_region_id" {
  description = "ID of the parent region"
  value       = netbox_region.parent.id
}

output "child_region_id" {
  description = "ID of the child region"
  value       = netbox_region.child.id
}

output "child_region_parent" {
  description = "Parent ID of the child region"
  value       = netbox_region.child.parent
}

output "grandchild_region_id" {
  description = "ID of the grandchild region"
  value       = netbox_region.grandchild.id
}

output "grandchild_region_parent" {
  description = "Parent ID of the grandchild region (should be child region)"
  value       = netbox_region.grandchild.parent
}
