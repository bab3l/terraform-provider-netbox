# Outputs to verify site group creation

output "root_id" {
  description = "ID of the root site group"
  value       = netbox_site_group.root.id
}

output "root_name" {
  description = "Name of the root site group"
  value       = netbox_site_group.root.name
}

output "root_slug" {
  description = "Slug of the root site group"
  value       = netbox_site_group.root.slug
}

output "child_id" {
  description = "ID of the child site group"
  value       = netbox_site_group.child.id
}

output "child_name" {
  description = "Name of the child site group"
  value       = netbox_site_group.child.name
}

output "child_parent" {
  description = "Parent ID of the child site group"
  value       = netbox_site_group.child.parent
}

output "grandchild_id" {
  description = "ID of the grandchild site group"
  value       = netbox_site_group.grandchild.id
}

output "grandchild_parent" {
  description = "Parent ID of the grandchild site group"
  value       = netbox_site_group.grandchild.parent
}

output "sibling_id" {
  description = "ID of the sibling site group"
  value       = netbox_site_group.sibling.id
}

# Verification: Child's parent should be root's ID
output "hierarchy_valid" {
  description = "Verification that hierarchy is correctly set"
  value       = netbox_site_group.child.parent == netbox_site_group.root.id && netbox_site_group.grandchild.parent == netbox_site_group.child.id
}
