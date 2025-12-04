output "prefix_basic_id" {
  description = "ID of the basic prefix"
  value       = netbox_prefix.basic.id
}

output "prefix_basic_prefix" {
  description = "Prefix of the basic prefix"
  value       = netbox_prefix.basic.prefix
}

output "prefix_complete_id" {
  description = "ID of the complete prefix"
  value       = netbox_prefix.complete.id
}

output "prefix_complete_status" {
  description = "Status of the complete prefix"
  value       = netbox_prefix.complete.status
}

output "prefix_complete_is_pool" {
  description = "Is pool flag of the complete prefix"
  value       = netbox_prefix.complete.is_pool
}

output "prefix_ipv6_id" {
  description = "ID of the IPv6 prefix"
  value       = netbox_prefix.ipv6.id
}

output "prefix_container_status" {
  description = "Status of the container prefix"
  value       = netbox_prefix.container.status
}

output "basic_prefix_valid" {
  description = "Validates basic prefix was created correctly"
  value       = netbox_prefix.basic.id != "" && netbox_prefix.basic.prefix == "10.0.0.0/24"
}

output "complete_prefix_valid" {
  description = "Validates complete prefix was created correctly"
  value       = netbox_prefix.complete.id != "" && netbox_prefix.complete.is_pool == true
}

output "vrf_association_valid" {
  description = "Validates prefix VRF association"
  value       = netbox_prefix.complete.vrf == netbox_vrf.test.id
}

output "site_association_valid" {
  description = "Validates prefix site association"
  value       = netbox_prefix.complete.site == netbox_site.test.id
}

output "tenant_association_valid" {
  description = "Validates prefix tenant association"
  value       = netbox_prefix.complete.tenant == netbox_tenant.test.id
}

output "vlan_association_valid" {
  description = "Validates prefix VLAN association"
  value       = netbox_prefix.complete.vlan == netbox_vlan.test.id
}

output "ipv6_prefix_valid" {
  description = "Validates IPv6 prefix was created correctly"
  value       = netbox_prefix.ipv6.id != "" && netbox_prefix.ipv6.prefix == "2001:db8::/32"
}
