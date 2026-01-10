# Lookup by ID
data "netbox_rack_type" "by_id" {
  id = "123"
}

# Lookup by slug
data "netbox_rack_type" "by_slug" {
  slug = "42u-2post-frame"
}

# Lookup by model name
data "netbox_rack_type" "by_model" {
  model = "FS-42"
}

# Use rack type data in other resources
output "rack_type_model" {
  value = data.netbox_rack_type.by_id.model
}

output "rack_type_slug" {
  value = data.netbox_rack_type.by_slug.slug
}

output "rack_type_manufacturer" {
  value = data.netbox_rack_type.by_model.manufacturer
}

output "rack_type_u_height" {
  value = data.netbox_rack_type.by_id.u_height
}

output "rack_type_form_factor" {
  value = data.netbox_rack_type.by_id.form_factor
}

# Access all custom fields
output "rack_type_custom_fields" {
  value       = data.netbox_rack_type.by_id.custom_fields
  description = "All custom fields defined in NetBox for this rack type"
}

# Access specific custom fields by name
output "rack_type_weight_capacity" {
  value       = try([for cf in data.netbox_rack_type.by_id.custom_fields : cf.value if cf.name == "weight_capacity_lbs"][0], null)
  description = "Example: accessing a numeric custom field"
}

output "rack_type_seismic_rated" {
  value       = try([for cf in data.netbox_rack_type.by_id.custom_fields : cf.value if cf.name == "seismic_rated"][0], null)
  description = "Example: accessing a boolean custom field"
}
