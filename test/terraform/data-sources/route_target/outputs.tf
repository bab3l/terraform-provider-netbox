# Route Target Data Source Outputs

# Look up by ID outputs
output "by_id_id" {
  value = data.netbox_route_target.by_id.id
}

output "by_id_name" {
  value = data.netbox_route_target.by_id.name
}

output "by_id_description" {
  value = data.netbox_route_target.by_id.description
}

# Look up by name outputs
output "by_name_id" {
  value = data.netbox_route_target.by_name.id
}

output "by_name_name" {
  value = data.netbox_route_target.by_name.name
}

output "by_name_description" {
  value = data.netbox_route_target.by_name.description
}

# Validation outputs
output "all_ids_match" {
  description = "Validates that all lookups return the same ID"
  value = alltrue([
    data.netbox_route_target.by_id.id == netbox_route_target.test.id,
    data.netbox_route_target.by_name.id == netbox_route_target.test.id
  ])
}

output "by_id_name_valid" {
  description = "Validates that lookup by ID returns correct name"
  value       = data.netbox_route_target.by_id.name == "65001:100"
}

output "by_name_name_valid" {
  description = "Validates that lookup by name returns correct name"
  value       = data.netbox_route_target.by_name.name == "65001:100"
}

output "descriptions_match" {
  description = "Validates that descriptions match the created resource"
  value       = data.netbox_route_target.by_id.description == "Test route target for data source"
}
