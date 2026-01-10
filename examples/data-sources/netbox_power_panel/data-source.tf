# Example 1: Lookup by ID
data "netbox_power_panel" "by_id" {
  id = "1"
}

# Example 2: Lookup by name only
data "netbox_power_panel" "by_name" {
  name = "Panel-A"
}

# Example 3: Lookup by name and site
data "netbox_power_panel" "by_name_and_site" {
  name = "Panel-B"
  site = "3"
}

# Use power panel data in other resources
output "power_panel_name" {
  value = data.netbox_power_panel.by_id.name
}

output "power_panel_site" {
  value = data.netbox_power_panel.by_id.site
}

output "power_panel_location" {
  value = data.netbox_power_panel.by_name.location
}

output "power_panel_description" {
  value = data.netbox_power_panel.by_name_and_site.description
}

# Access all custom fields
output "power_panel_custom_fields" {
  value       = data.netbox_power_panel.by_id.custom_fields
  description = "All custom fields defined in NetBox for this power panel"
}

# Access specific custom fields by name
output "power_panel_capacity" {
  value       = try([for cf in data.netbox_power_panel.by_id.custom_fields : cf.value if cf.name == "capacity_amps"][0], null)
  description = "Example: accessing a numeric custom field"
}

output "power_panel_utility_company" {
  value       = try([for cf in data.netbox_power_panel.by_id.custom_fields : cf.value if cf.name == "utility_company"][0], null)
  description = "Example: accessing a text custom field"
}
