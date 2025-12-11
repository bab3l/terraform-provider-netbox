output "basic_id" {
  value = netbox_event_rule.basic.id
}

output "basic_name" {
  value = netbox_event_rule.basic.name
}

output "multiple_types_id" {
  value = netbox_event_rule.multiple_types.id
}

output "multiple_types_object_types" {
  value = netbox_event_rule.multiple_types.object_types
}

output "multiple_types_event_types" {
  value = netbox_event_rule.multiple_types.event_types
}

output "complete_id" {
  value = netbox_event_rule.complete.id
}

output "complete_description" {
  value = netbox_event_rule.complete.description
}

output "complete_tags" {
  value = netbox_event_rule.complete.tags
}

output "disabled_id" {
  value = netbox_event_rule.disabled.id
}

output "disabled_enabled" {
  value = netbox_event_rule.disabled.enabled
}
