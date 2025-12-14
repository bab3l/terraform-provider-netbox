terraform {
  required_version = ">= 1.0"
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

# Create a choice set to look up
resource "netbox_custom_field_choice_set" "test" {
  name        = "test_choice_set_lookup"
  description = "Test choice set for data source lookup"
  extra_choices = [
    { value = "opt1", label = "Option 1" },
    { value = "opt2", label = "Option 2" },
    { value = "opt3", label = "Option 3" },
  ]
}

# Look up by ID
data "netbox_custom_field_choice_set" "by_id" {
  id = netbox_custom_field_choice_set.test.id
}

# Look up by name
data "netbox_custom_field_choice_set" "by_name" {
  name = netbox_custom_field_choice_set.test.name
}

output "lookup_by_id_name" {
  value = data.netbox_custom_field_choice_set.by_id.name
}

output "lookup_by_name_id" {
  value = data.netbox_custom_field_choice_set.by_name.id
}

output "lookup_by_id_description" {
  value = data.netbox_custom_field_choice_set.by_id.description
}

output "lookup_by_name_extra_choices_count" {
  value = length(data.netbox_custom_field_choice_set.by_name.extra_choices)
}
