# Look up aggregate by prefix
data "netbox_aggregate" "by_prefix" {
  prefix = "10.0.0.0/8"
}

# Look up aggregate by ID
data "netbox_aggregate" "by_id" {
  id = "123"
}

# Use aggregate data in other resources
output "aggregate_prefix" {
  value = data.netbox_aggregate.by_id.prefix
}

output "aggregate_rir" {
  value = data.netbox_aggregate.by_prefix.rir
}

output "aggregate_tenant" {
  value = data.netbox_aggregate.by_id.tenant
}

output "aggregate_description" {
  value = data.netbox_aggregate.by_id.description
}

# Access all custom fields
output "aggregate_custom_fields" {
  value       = data.netbox_aggregate.by_id.custom_fields
  description = "All custom fields defined in NetBox for this aggregate"
}

# Access specific custom fields by name
output "aggregate_allocation_date" {
  value       = try([for cf in data.netbox_aggregate.by_id.custom_fields : cf.value if cf.name == "allocation_date"][0], null)
  description = "Example: accessing a date custom field"
}

output "aggregate_ipv6_enabled" {
  value       = try([for cf in data.netbox_aggregate.by_id.custom_fields : cf.value if cf.name == "ipv6_enabled"][0], null)
  description = "Example: accessing a boolean custom field"
}
