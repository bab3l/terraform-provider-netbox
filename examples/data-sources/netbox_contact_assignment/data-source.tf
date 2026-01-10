# Look up contact assignment by ID
data "netbox_contact_assignment" "by_id" {
  id = "123"
}

# Use contact assignment data in other resources
output "assignment_contact" {
  value = data.netbox_contact_assignment.by_id.contact
}

output "assignment_object_type" {
  value = data.netbox_contact_assignment.by_id.object_type
}

output "assignment_object_id" {
  value = data.netbox_contact_assignment.by_id.object_id
}

output "assignment_role" {
  value = data.netbox_contact_assignment.by_id.role
}

output "assignment_priority" {
  value = data.netbox_contact_assignment.by_id.priority
}

# Access all custom fields
output "assignment_custom_fields" {
  value       = data.netbox_contact_assignment.by_id.custom_fields
  description = "All custom fields defined in NetBox for this contact assignment"
}

# Access specific custom fields by name
output "assignment_notification_method" {
  value       = try([for cf in data.netbox_contact_assignment.by_id.custom_fields : cf.value if cf.name == "notification_method"][0], null)
  description = "Example: accessing a select custom field"
}

output "assignment_escalation_level" {
  value       = try([for cf in data.netbox_contact_assignment.by_id.custom_fields : cf.value if cf.name == "escalation_level"][0], null)
  description = "Example: accessing a numeric custom field"
}
