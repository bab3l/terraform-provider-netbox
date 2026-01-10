data "netbox_circuit_type" "by_id" {
  id = "123"
}

data "netbox_circuit_type" "by_name" {
  name = "Internet Transit"
}

data "netbox_circuit_type" "by_slug" {
  slug = "internet-transit"
}

output "by_id" {
  value = data.netbox_circuit_type.by_id.name
}

output "by_name" {
  value = data.netbox_circuit_type.by_name.id
}

output "by_slug" {
  value = data.netbox_circuit_type.by_slug.id
}

output "circuit_type_slug" {
  value = data.netbox_circuit_type.by_id.slug
}

output "circuit_type_description" {
  value = data.netbox_circuit_type.by_id.description
}

# Access all custom fields
output "circuit_type_custom_fields" {
  value       = data.netbox_circuit_type.by_id.custom_fields
  description = "All custom fields defined in NetBox for this circuit type"
}

# Access specific custom field by name
output "circuit_type_sla_target" {
  value       = try([for cf in data.netbox_circuit_type.by_id.custom_fields : cf.value if cf.name == "sla_target_uptime"][0], null)
  description = "Example: accessing a numeric custom field for SLA target"
}

output "circuit_type_requires_redundancy" {
  value       = try([for cf in data.netbox_circuit_type.by_id.custom_fields : cf.value if cf.name == "requires_redundancy"][0], null)
  description = "Example: accessing a boolean custom field for redundancy requirement"
}
