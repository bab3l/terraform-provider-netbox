# Look up by ID
data "netbox_custom_field_choice_set" "by_id" {
  id = "123"
}

# Look up by name
data "netbox_custom_field_choice_set" "by_name" {
  name = "priority_levels"
}

output "custom_field_choice_set_by_id" {
  value = data.netbox_custom_field_choice_set.by_id
}

output "custom_field_choice_set_by_name" {
  value = data.netbox_custom_field_choice_set.by_name
}
