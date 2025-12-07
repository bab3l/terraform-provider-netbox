# Virtual Chassis Outputs

# Basic virtual chassis outputs
output "basic_id" {
  value = netbox_virtual_chassis.basic.id
}

output "basic_name" {
  value = netbox_virtual_chassis.basic.name
}

# Complete virtual chassis outputs
output "complete_id" {
  value = netbox_virtual_chassis.complete.id
}

output "complete_name" {
  value = netbox_virtual_chassis.complete.name
}

output "complete_domain" {
  value = netbox_virtual_chassis.complete.domain
}
