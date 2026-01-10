# Look up a custom field choice set by ID
data "netbox_custom_field_choice_set" "by_id" {
  id = "123"
}

# Look up a custom field choice set by name
data "netbox_custom_field_choice_set" "by_name" {
  name = "priority_levels"
}

# Individual attribute outputs
output "choice_set_id" {
  value       = data.netbox_custom_field_choice_set.by_id.id
  description = "The unique ID of the custom field choice set"
}

output "choice_set_name" {
  value       = data.netbox_custom_field_choice_set.by_name.name
  description = "The name of the custom field choice set"
}

output "choice_set_choices" {
  value       = data.netbox_custom_field_choice_set.by_name.choices
  description = "List of available choices in this custom field choice set"
}

output "choice_set_description" {
  value       = data.netbox_custom_field_choice_set.by_name.description
  description = "Description of the custom field choice set"
}

# Access all custom fields
output "choice_set_custom_fields" {
  value       = data.netbox_custom_field_choice_set.by_id.custom_fields
  description = "All custom fields defined in NetBox for this choice set"
}

# Access specific custom field by name
output "choice_set_custom_field_example" {
  value       = try([for cf in data.netbox_custom_field_choice_set.by_id.custom_fields : cf.value if cf.name == "internal_id"][0], null)
  description = "Example: accessing a specific custom field value (internal_id)"
}
