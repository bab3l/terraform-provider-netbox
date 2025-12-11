terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

# Basic custom field choice set with extra choices
resource "netbox_custom_field_choice_set" "priority_levels" {
  name = "priority_levels"
  extra_choices = [
    { value = "critical", label = "Critical" },
    { value = "high",     label = "High" },
    { value = "medium",   label = "Medium" },
    { value = "low",      label = "Low" },
  ]
}

# Custom field choice set with full options (no alphabetical sorting to avoid order mismatch)
resource "netbox_custom_field_choice_set" "full_example" {
  name                 = "service_types"
  description          = "Types of services provided"
  order_alphabetically = false
  extra_choices = [
    { value = "web",      label = "Web Services" },
    { value = "database", label = "Database Services" },
    { value = "cache",    label = "Caching Services" },
    { value = "queue",    label = "Message Queue" },
  ]
}

# Custom field choice set with base choices (country codes)
resource "netbox_custom_field_choice_set" "country_extended" {
  name         = "country_extended"
  base_choices = "ISO_3166"
  description  = "Country codes with additional custom options"
  extra_choices = [
    { value = "EMEA",   label = "Europe, Middle East & Africa" },
    { value = "APAC",   label = "Asia Pacific" },
    { value = "AMER",   label = "Americas" },
  ]
}

output "priority_id" {
  value = netbox_custom_field_choice_set.priority_levels.id
}

output "full_example_id" {
  value = netbox_custom_field_choice_set.full_example.id
}

output "country_extended_id" {
  value = netbox_custom_field_choice_set.country_extended.id
}
