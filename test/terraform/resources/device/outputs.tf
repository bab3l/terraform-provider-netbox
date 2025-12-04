output "device_basic_id" {
  description = "ID of the basic device"
  value       = netbox_device.basic.id
}

output "device_basic_name" {
  description = "Name of the basic device"
  value       = netbox_device.basic.name
}

output "device_complete_id" {
  description = "ID of the complete device"
  value       = netbox_device.complete.id
}

output "device_complete_serial" {
  description = "Serial number of the complete device"
  value       = netbox_device.complete.serial
}

output "device_complete_asset_tag" {
  description = "Asset tag of the complete device"
  value       = netbox_device.complete.asset_tag
}

output "device_racked_id" {
  description = "ID of the racked device"
  value       = netbox_device.racked.id
}

output "device_racked_position" {
  description = "Position of the racked device"
  value       = netbox_device.racked.position
}

output "device_planned_status" {
  description = "Status of the planned device"
  value       = netbox_device.planned.status
}

output "basic_device_valid" {
  description = "Validates basic device was created correctly"
  value       = netbox_device.basic.id != "" && netbox_device.basic.name == "Basic Test Device"
}

output "complete_device_valid" {
  description = "Validates complete device was created correctly"
  value       = netbox_device.complete.id != "" && netbox_device.complete.serial == "SN-COMPLETE-001"
}

output "site_reference_valid" {
  description = "Validates device site reference"
  value       = netbox_device.basic.site == netbox_site.test.id
}

output "device_type_reference_valid" {
  description = "Validates device type reference"
  value       = netbox_device.basic.device_type == netbox_device_type.test.id
}

output "role_reference_valid" {
  description = "Validates device role reference"
  value       = netbox_device.basic.role == netbox_device_role.test.id
}

output "tenant_association_valid" {
  description = "Validates device tenant association"
  value       = netbox_device.complete.tenant == netbox_tenant.test.id
}

output "platform_association_valid" {
  description = "Validates device platform association"
  value       = netbox_device.complete.platform == netbox_platform.test.id
}

output "rack_association_valid" {
  description = "Validates device rack association"
  value       = netbox_device.racked.rack == netbox_rack.test.id
}
