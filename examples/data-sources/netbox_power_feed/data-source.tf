# Example 1: Lookup by ID
data "netbox_power_feed" "by_id" {
  id = "1"
}

# Example 2: Lookup by power_panel and name
data "netbox_power_feed" "by_panel_and_name" {
  power_panel = "5"
  name        = "Feed-A"
}

# Use power feed data in other resources
output "power_feed_name" {
  value = data.netbox_power_feed.by_id.name
}

output "power_feed_status" {
  value = data.netbox_power_feed.by_panel_and_name.status
}

output "power_feed_type" {
  value = data.netbox_power_feed.by_id.type
}

output "power_feed_voltage" {
  value = data.netbox_power_feed.by_id.voltage
}

output "power_feed_amperage" {
  value = data.netbox_power_feed.by_id.amperage
}

# Access all custom fields
output "power_feed_custom_fields" {
  value       = data.netbox_power_feed.by_id.custom_fields
  description = "All custom fields defined in NetBox for this power feed"
}

# Access specific custom fields by name
output "power_feed_circuit_id" {
  value       = try([for cf in data.netbox_power_feed.by_id.custom_fields : cf.value if cf.name == "circuit_id"][0], null)
  description = "Example: accessing a text custom field"
}

output "power_feed_redundant" {
  value       = try([for cf in data.netbox_power_feed.by_id.custom_fields : cf.value if cf.name == "is_redundant"][0], null)
  description = "Example: accessing a boolean custom field"
}
