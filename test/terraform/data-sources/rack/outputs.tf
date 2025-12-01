# Outputs to verify rack data source lookups

# Outputs for lookup by ID
output "by_id_rack_id" {
  description = "ID of the rack looked up by ID"
  value       = data.netbox_rack.by_id.id
}

output "by_id_rack_name" {
  description = "Name of the rack looked up by ID"
  value       = data.netbox_rack.by_id.name
}

output "by_id_rack_site" {
  description = "Site ID of the rack looked up by ID"
  value       = data.netbox_rack.by_id.site
}

output "by_id_rack_location" {
  description = "Location ID of the rack looked up by ID"
  value       = data.netbox_rack.by_id.location
}

output "by_id_rack_tenant" {
  description = "Tenant ID of the rack looked up by ID"
  value       = data.netbox_rack.by_id.tenant
}

output "by_id_rack_status" {
  description = "Status of the rack looked up by ID"
  value       = data.netbox_rack.by_id.status
}

output "by_id_rack_serial" {
  description = "Serial number of the rack looked up by ID"
  value       = data.netbox_rack.by_id.serial
}

output "by_id_rack_asset_tag" {
  description = "Asset tag of the rack looked up by ID"
  value       = data.netbox_rack.by_id.asset_tag
}

output "by_id_rack_u_height" {
  description = "U height of the rack looked up by ID"
  value       = data.netbox_rack.by_id.u_height
}

output "by_id_rack_description" {
  description = "Description of the rack looked up by ID"
  value       = data.netbox_rack.by_id.description
}

output "by_id_rack_width" {
  description = "Width of the rack looked up by ID"
  value       = data.netbox_rack.by_id.width
}

# Outputs for lookup by name
output "by_name_rack_id" {
  description = "ID of the rack looked up by name"
  value       = data.netbox_rack.by_name.id
}

output "by_name_rack_name" {
  description = "Name of the rack looked up by name"
  value       = data.netbox_rack.by_name.name
}

output "by_name_rack_site" {
  description = "Site ID of the rack looked up by name"
  value       = data.netbox_rack.by_name.site
}

output "by_name_rack_status" {
  description = "Status of the rack looked up by name"
  value       = data.netbox_rack.by_name.status
}

# Outputs for reserved rack
output "reserved_rack_id" {
  description = "ID of the reserved rack"
  value       = data.netbox_rack.reserved_by_id.id
}

output "reserved_rack_status" {
  description = "Status of the reserved rack"
  value       = data.netbox_rack.reserved_by_id.status
}

output "reserved_rack_u_height" {
  description = "U height of the reserved rack"
  value       = data.netbox_rack.reserved_by_id.u_height
}

# Verify consistency between lookup methods
output "id_matches_between_lookups" {
  description = "Verify ID lookup matches name lookup"
  value       = data.netbox_rack.by_id.id == data.netbox_rack.by_name.id
}
