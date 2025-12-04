output "vlan_basic_id" {
  description = "ID of the basic VLAN"
  value       = netbox_vlan.basic.id
}

output "vlan_basic_vid" {
  description = "VID of the basic VLAN"
  value       = netbox_vlan.basic.vid
}

output "vlan_basic_name" {
  description = "Name of the basic VLAN"
  value       = netbox_vlan.basic.name
}

output "vlan_complete_id" {
  description = "ID of the complete VLAN"
  value       = netbox_vlan.complete.id
}

output "vlan_complete_status" {
  description = "Status of the complete VLAN"
  value       = netbox_vlan.complete.status
}

output "vlan_reserved_status" {
  description = "Status of the reserved VLAN"
  value       = netbox_vlan.reserved.status
}

output "basic_vlan_valid" {
  description = "Validates basic VLAN was created correctly"
  value       = netbox_vlan.basic.id != "" && netbox_vlan.basic.vid == 100
}

output "complete_vlan_valid" {
  description = "Validates complete VLAN was created correctly"
  value       = netbox_vlan.complete.id != "" && netbox_vlan.complete.vid == 200 && netbox_vlan.complete.status == "active"
}

output "vlan_group_association_valid" {
  description = "Validates VLAN group association"
  value       = netbox_vlan.complete.group == netbox_vlan_group.test.id
}

output "vlan_site_association_valid" {
  description = "Validates VLAN site association"
  value       = netbox_vlan.with_site.site == netbox_site.test.id
}

output "vlan_tenant_association_valid" {
  description = "Validates VLAN tenant association"
  value       = netbox_vlan.complete.tenant == netbox_tenant.test.id
}
