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

output "event_rule_description" {
  value = data.netbox_event_rule.device_changes.description
}

output "event_rule_type_create" {
  value = data.netbox_event_rule.device_changes.type_create
}

output "event_rule_type_update" {
  value = data.netbox_event_rule.device_changes.type_update
}

output "event_rule_type_delete" {
  value = data.netbox_event_rule.device_changes.type_delete
}

# Note: Event rules do not support custom fields in NetBox API
output "event_rule_triggers" {
  value       = data.netbox_event_rule.device_changes.triggers
  description = "Event rule triggers configuration"
}

output "event_rule_actions" {
  value       = data.netbox_event_rule.device_changes.actions
  description = "Event rule actions configuration"
}
