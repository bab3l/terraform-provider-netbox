output "event_rule_id" {
  value = data.netbox_event_rule.test.id
}

output "event_rule_name" {
  value = data.netbox_event_rule.test.name
}

output "event_rule_object_types" {
  value = data.netbox_event_rule.test.object_types
}

output "event_rule_event_types" {
  value = data.netbox_event_rule.test.event_types
}

output "event_rule_action_type" {
  value = data.netbox_event_rule.test.action_type
}

output "event_rule_enabled" {
  value = data.netbox_event_rule.test.enabled
}

output "event_rule_description" {
  value = data.netbox_event_rule.test.description
}
