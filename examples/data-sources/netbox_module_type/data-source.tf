# Example 1: Lookup by ID
data "netbox_module_type" "by_id" {
  id = 1
}

output "module_type_by_id" {
  value       = data.netbox_module_type.by_id.model
  description = "Module type model when looked up by ID"
}

# Example 2: Lookup by model only
data "netbox_module_type" "by_model" {
  model = "XM-100"
}

output "module_type_by_model" {
  value       = data.netbox_module_type.by_model.manufacturer
  description = "Manufacturer of module type when looked up by model"
}

# Example 3: Lookup by model and manufacturer_id
data "netbox_module_type" "by_model_and_manufacturer" {
  model           = "XM-200"
  manufacturer_id = 10
}

output "module_type_by_model_and_manufacturer" {
  value       = data.netbox_module_type.by_model_and_manufacturer.part_number
  description = "Part number of module type when looked up by model and manufacturer"
}
