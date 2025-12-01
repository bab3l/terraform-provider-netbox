# Outputs to verify rack resource creation

output "basic_rack_id" {
  description = "ID of the basic rack"
  value       = netbox_rack.basic.id
}

output "basic_rack_name" {
  description = "Name of the basic rack"
  value       = netbox_rack.basic.name
}

output "basic_rack_site" {
  description = "Site ID of the basic rack"
  value       = netbox_rack.basic.site
}

output "basic_rack_status" {
  description = "Status of the basic rack"
  value       = netbox_rack.basic.status
}

output "complete_rack_id" {
  description = "ID of the complete rack"
  value       = netbox_rack.complete.id
}

output "complete_rack_serial" {
  description = "Serial number of the complete rack"
  value       = netbox_rack.complete.serial
}

output "complete_rack_asset_tag" {
  description = "Asset tag of the complete rack"
  value       = netbox_rack.complete.asset_tag
}

output "complete_rack_u_height" {
  description = "U height of the complete rack"
  value       = netbox_rack.complete.u_height
}

output "complete_rack_outer_width" {
  description = "Outer width of the complete rack"
  value       = netbox_rack.complete.outer_width
}

output "complete_rack_outer_depth" {
  description = "Outer depth of the complete rack"
  value       = netbox_rack.complete.outer_depth
}

output "complete_rack_description" {
  description = "Description of the complete rack"
  value       = netbox_rack.complete.description
}

output "complete_rack_location" {
  description = "Location ID of the complete rack"
  value       = netbox_rack.complete.location
}

output "complete_rack_tenant" {
  description = "Tenant ID of the complete rack"
  value       = netbox_rack.complete.tenant
}

output "descending_rack_desc_units" {
  description = "Whether the descending rack uses descending units"
  value       = netbox_rack.descending.desc_units
}

output "reserved_rack_status" {
  description = "Status of the reserved rack"
  value       = netbox_rack.reserved.status
}

output "planned_rack_status" {
  description = "Status of the planned rack"
  value       = netbox_rack.planned.status
}

output "deprecated_rack_status" {
  description = "Status of the deprecated rack"
  value       = netbox_rack.deprecated.status
}

output "imperial_rack_outer_unit" {
  description = "Outer unit of the imperial rack"
  value       = netbox_rack.imperial.outer_unit
}

output "weight_lb_rack_weight_unit" {
  description = "Weight unit of the weight_lb rack"
  value       = netbox_rack.weight_lb.weight_unit
}

output "small_rack_u_height" {
  description = "U height of the small rack"
  value       = netbox_rack.small.u_height
}

output "large_rack_u_height" {
  description = "U height of the large rack"
  value       = netbox_rack.large.u_height
}

output "complete_rack_width" {
  description = "Width of the complete rack (default computed)"
  value       = netbox_rack.complete.width
}
