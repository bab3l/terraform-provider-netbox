# Look up by ID
data "netbox_custom_field_choice_set" "by_id" {
  id = "123"
}

# Look up by name
data "netbox_custom_field_choice_set" "by_name" {
  name = "priority_levels"
}
