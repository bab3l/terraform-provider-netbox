# Outputs to verify tenant group creation

output "root_id" {
  description = "ID of the root tenant group"
  value       = netbox_tenant_group.root.id
}

output "root_name" {
  description = "Name of the root tenant group"
  value       = netbox_tenant_group.root.name
}

output "root_slug" {
  description = "Slug of the root tenant group"
  value       = netbox_tenant_group.root.slug
}

output "child_id" {
  description = "ID of the child tenant group"
  value       = netbox_tenant_group.child.id
}

output "child_name" {
  description = "Name of the child tenant group"
  value       = netbox_tenant_group.child.name
}

output "child_parent" {
  description = "Parent ID of the child tenant group"
  value       = netbox_tenant_group.child.parent
}

output "grandchild_id" {
  description = "ID of the grandchild tenant group"
  value       = netbox_tenant_group.grandchild.id
}

output "grandchild_parent" {
  description = "Parent ID of the grandchild tenant group"
  value       = netbox_tenant_group.grandchild.parent
}

output "sibling_id" {
  description = "ID of the sibling tenant group"
  value       = netbox_tenant_group.sibling.id
}

output "with_tenant_id" {
  description = "ID of the tenant group with tenant"
  value       = netbox_tenant_group.with_tenant.id
}

output "tenant_in_group_id" {
  description = "ID of the tenant in the group"
  value       = netbox_tenant.in_group.id
}

output "tenant_in_group_group" {
  description = "Group ID of the tenant"
  value       = netbox_tenant.in_group.group
}

# Verification
output "hierarchy_valid" {
  description = "Verification that hierarchy is correctly set"
  value       = netbox_tenant_group.child.parent == netbox_tenant_group.root.id && netbox_tenant_group.grandchild.parent == netbox_tenant_group.child.id
}

output "tenant_group_assignment_valid" {
  description = "Verification that tenant is correctly assigned to group"
  value       = netbox_tenant.in_group.group == netbox_tenant_group.with_tenant.id
}
