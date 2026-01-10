# Example 1: Lookup by ID
data "netbox_module_type" "by_id" {
  id = "1"
}

# Example 2: Lookup by model only
data "netbox_module_type" "by_model" {
  model = "XM-100"
}

# Example 3: Lookup by model and manufacturer_id
data "netbox_module_type" "by_model_and_manufacturer" {
  model           = "XM-200"
  manufacturer_id = "10"
}

# Use module type data in other resources
output "module_type_model" {
  value = data.netbox_module_type.by_id.model
}

output "module_type_manufacturer" {
  value = data.netbox_module_type.by_model.manufacturer
}

output "module_type_part_number" {
  value = data.netbox_module_type.by_model_and_manufacturer.part_number
}

output "module_type_weight" {
  value = data.netbox_module_type.by_id.weight
}

# Access all custom fields
output "module_type_custom_fields" {
  value       = data.netbox_module_type.by_id.custom_fields
  description = "All custom fields defined in NetBox for this module type"
}

# Access specific custom fields by name
output "module_type_power_consumption" {
  value       = try([for cf in data.netbox_module_type.by_id.custom_fields : cf.value if cf.name == "power_consumption_watts"][0], null)
  description = "Example: accessing a numeric custom field"
}

output "module_type_compatible_chassis" {
  value       = try([for cf in data.netbox_module_type.by_id.custom_fields : cf.value if cf.name == "compatible_chassis"][0], null)
  description = "Example: accessing a multiselect custom field"
}
