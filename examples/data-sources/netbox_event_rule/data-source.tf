# Look up an existing event rule by ID
data "netbox_event_rule" "device_changes" {
  id = "1"
}

# Use the event rule data
output "event_rule_name" {
  value = data.netbox_event_rule.device_changes.name
}

output "event_rule_object_types" {
  value = data.netbox_event_rule.device_changes.object_types
}

output "event_rule_enabled" {
  value = data.netbox_event_rule.device_changes.enabled
}
