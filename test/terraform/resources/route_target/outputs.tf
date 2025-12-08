# Route Target Outputs

# Basic route target outputs
output "basic_id" {
  value = netbox_route_target.basic.id
}

output "basic_name" {
  value = netbox_route_target.basic.name
}

output "basic_id_valid" {
  description = "Validates that the basic route target was created with an ID"
  value       = netbox_route_target.basic.id != null && netbox_route_target.basic.id != ""
}

output "basic_name_valid" {
  description = "Validates that the basic route target name matches the input"
  value       = netbox_route_target.basic.name == "65000:100"
}

# Complete route target outputs
output "complete_id" {
  value = netbox_route_target.complete.id
}

output "complete_name" {
  value = netbox_route_target.complete.name
}

output "complete_description" {
  value = netbox_route_target.complete.description
}

output "complete_comments" {
  value = netbox_route_target.complete.comments
}

output "complete_name_valid" {
  description = "Validates that the complete route target name matches the input"
  value       = netbox_route_target.complete.name == "65000:200"
}

output "complete_description_valid" {
  description = "Validates that the description was set correctly"
  value       = netbox_route_target.complete.description == "Test route target for VRF export"
}

# Route target with tenant outputs
output "with_tenant_id" {
  value = netbox_route_target.with_tenant.id
}

output "with_tenant_name" {
  value = netbox_route_target.with_tenant.name
}

output "with_tenant_tenant" {
  value = netbox_route_target.with_tenant.tenant
}

output "with_tenant_tenant_valid" {
  description = "Validates that the tenant was associated correctly"
  value       = netbox_route_target.with_tenant.tenant == netbox_tenant.rt_test.id
}

# Aggregate validation output
output "all_tests_passed" {
  description = "Validates all route target tests passed"
  value = alltrue([
    netbox_route_target.basic.id != null && netbox_route_target.basic.id != "",
    netbox_route_target.basic.name == "65000:100",
    netbox_route_target.complete.name == "65000:200",
    netbox_route_target.complete.description == "Test route target for VRF export",
    netbox_route_target.with_tenant.tenant == netbox_tenant.rt_test.id
  ])
}
