output "vrf_basic_id" {
  description = "ID of the basic VRF"
  value       = netbox_vrf.basic.id
}

output "vrf_basic_name" {
  description = "Name of the basic VRF"
  value       = netbox_vrf.basic.name
}

output "vrf_complete_id" {
  description = "ID of the complete VRF"
  value       = netbox_vrf.complete.id
}

output "vrf_complete_rd" {
  description = "Route distinguisher of the complete VRF"
  value       = netbox_vrf.complete.rd
}

output "vrf_with_tenant_id" {
  description = "ID of the VRF with tenant"
  value       = netbox_vrf.with_tenant.id
}

output "basic_vrf_valid" {
  description = "Validates basic VRF was created correctly"
  value       = netbox_vrf.basic.id != "" && netbox_vrf.basic.name == "Basic Test VRF"
}

output "complete_vrf_valid" {
  description = "Validates complete VRF was created correctly"
  value       = netbox_vrf.complete.id != "" && netbox_vrf.complete.rd == "65000:100"
}

output "tenant_association_valid" {
  description = "Validates VRF tenant association was created correctly"
  value       = netbox_vrf.with_tenant.tenant == netbox_tenant.test.id
}
