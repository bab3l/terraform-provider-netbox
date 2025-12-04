output "ip_address_basic_id" {
  description = "ID of the basic IP address"
  value       = netbox_ip_address.basic.id
}

output "ip_address_basic_address" {
  description = "Address of the basic IP address"
  value       = netbox_ip_address.basic.address
}

output "ip_address_complete_id" {
  description = "ID of the complete IP address"
  value       = netbox_ip_address.complete.id
}

output "ip_address_complete_dns_name" {
  description = "DNS name of the complete IP address"
  value       = netbox_ip_address.complete.dns_name
}

output "ip_address_complete_status" {
  description = "Status of the complete IP address"
  value       = netbox_ip_address.complete.status
}

output "ip_address_ipv6_id" {
  description = "ID of the IPv6 address"
  value       = netbox_ip_address.ipv6.id
}

output "ip_address_reserved_status" {
  description = "Status of the reserved IP address"
  value       = netbox_ip_address.reserved.status
}

output "ip_address_dhcp_status" {
  description = "Status of the DHCP IP address"
  value       = netbox_ip_address.dhcp.status
}

output "basic_ip_address_valid" {
  description = "Validates basic IP address was created correctly"
  value       = netbox_ip_address.basic.id != "" && netbox_ip_address.basic.address == "10.100.0.1/24"
}

output "complete_ip_address_valid" {
  description = "Validates complete IP address was created correctly"
  value       = netbox_ip_address.complete.id != "" && netbox_ip_address.complete.dns_name == "test-server.example.com"
}

output "vrf_association_valid" {
  description = "Validates IP address VRF association"
  value       = netbox_ip_address.complete.vrf == netbox_vrf.test.id
}

output "tenant_association_valid" {
  description = "Validates IP address tenant association"
  value       = netbox_ip_address.complete.tenant == netbox_tenant.test.id
}

output "ipv6_address_valid" {
  description = "Validates IPv6 address was created correctly"
  value       = netbox_ip_address.ipv6.id != "" && netbox_ip_address.ipv6.address == "2001:db8::1/64"
}
