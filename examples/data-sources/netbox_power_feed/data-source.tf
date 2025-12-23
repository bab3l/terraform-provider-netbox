# Example 1: Lookup by ID
data "netbox_power_feed" "by_id" {
  id = 1
}

output "power_feed_by_id" {
  value       = data.netbox_power_feed.by_id.name
  description = "Power feed name when looked up by ID"
}

# Example 2: Lookup by power_panel and name
data "netbox_power_feed" "by_panel_and_name" {
  power_panel = 5
  name        = "Feed-A"
}

output "power_feed_by_panel_and_name" {
  value       = data.netbox_power_feed.by_panel_and_name.status
  description = "Power feed status when looked up by panel and name"
}
