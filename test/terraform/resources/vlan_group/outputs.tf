output "vlan_group_basic_id" {
  description = "ID of the basic VLAN group"
  value       = netbox_vlan_group.basic.id
}

output "vlan_group_basic_name" {
  description = "Name of the basic VLAN group"
  value       = netbox_vlan_group.basic.name
}

output "vlan_group_basic_slug" {
  description = "Slug of the basic VLAN group"
  value       = netbox_vlan_group.basic.slug
}

output "vlan_group_complete_id" {
  description = "ID of the complete VLAN group"
  value       = netbox_vlan_group.complete.id
}

output "vlan_group_complete_description" {
  description = "Description of the complete VLAN group"
  value       = netbox_vlan_group.complete.description
}

output "basic_vlan_group_valid" {
  description = "Validates basic VLAN group was created correctly"
  value       = netbox_vlan_group.basic.id != "" && netbox_vlan_group.basic.slug == "basic-vlan-group"
}

output "complete_vlan_group_valid" {
  description = "Validates complete VLAN group was created correctly"
  value       = netbox_vlan_group.complete.id != "" && netbox_vlan_group.complete.description == "Complete VLAN group for integration testing"
}

output "site_scoped_vlan_group_valid" {
  description = "Validates site-scoped VLAN group was created correctly"
  value       = netbox_vlan_group.site_scoped.id != "" && netbox_vlan_group.site_scoped.scope_type == "dcim.site"
}
