# Example 1: Lookup by ID
data "netbox_power_panel" "by_id" {
  id = 1
}

output "power_panel_by_id" {
  value       = data.netbox_power_panel.by_id.name
  description = "Power panel name when looked up by ID"
}

# Example 2: Lookup by name only
data "netbox_power_panel" "by_name" {
  name = "Panel-A"
}

output "power_panel_by_name" {
  value       = data.netbox_power_panel.by_name.description
  description = "Power panel description"
}

# Example 3: Lookup by name and site
data "netbox_power_panel" "by_name_and_site" {
  name = "Panel-B"
  site = 3
}

output "power_panel_by_name_and_site" {
  value       = data.netbox_power_panel.by_name_and_site.display_name
  description = "Power panel display name when looked up by name and site"
}
