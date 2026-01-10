# Look up a virtual chassis by ID
data "netbox_virtual_chassis" "by_id" {
  id = "1"
}

# Look up a virtual chassis by name
data "netbox_virtual_chassis" "by_name" {
  name = "test-virtual-chassis"
}

# Use virtual chassis data in outputs
output "by_id" {
  value = data.netbox_virtual_chassis.by_id.name
}

output "by_name" {
  value = data.netbox_virtual_chassis.by_name.id
}

output "master_device" {
  value = data.netbox_virtual_chassis.by_name.master
}

output "chassis_domain" {
  value = data.netbox_virtual_chassis.by_id.domain
}

output "chassis_description" {
  value = data.netbox_virtual_chassis.by_id.description
}

# Access all custom fields
output "chassis_custom_fields" {
  value       = data.netbox_virtual_chassis.by_id.custom_fields
  description = "All custom fields defined in NetBox for this virtual chassis"
}

# Access specific custom field by name
output "chassis_stack_id" {
  value       = try([for cf in data.netbox_virtual_chassis.by_id.custom_fields : cf.value if cf.name == "stack_id"][0], null)
  description = "Example: accessing a text custom field for stack ID"
}

output "chassis_member_count" {
  value       = try([for cf in data.netbox_virtual_chassis.by_id.custom_fields : cf.value if cf.name == "member_count"][0], null)
  description = "Example: accessing a numeric custom field for member count"
}
