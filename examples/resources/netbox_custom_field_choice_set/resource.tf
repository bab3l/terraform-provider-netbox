resource "netbox_custom_field_choice_set" "example" {
  name = "priority_levels"
  extra_choices = [
    { value = "critical", label = "Critical" },
    { value = "high", label = "High" },
    { value = "medium", label = "Medium" },
    { value = "low", label = "Low" },
  ]
}
