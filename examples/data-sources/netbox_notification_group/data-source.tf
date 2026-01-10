# Look up a notification group by ID
data "netbox_notification_group" "by_id" {
  id = "1"
}

# Look up a notification group by name
data "netbox_notification_group" "by_name" {
  name = "Administrators"
}

# Individual attribute outputs
output "notification_group_id" {
  value       = data.netbox_notification_group.by_id.id
  description = "The unique ID of the notification group"
}

output "notification_group_name" {
  value       = data.netbox_notification_group.by_name.name
  description = "The name of the notification group"
}

output "notification_group_description" {
  value       = data.netbox_notification_group.by_name.description
  description = "Description of the notification group"
}

output "notification_group_group_ids" {
  value       = data.netbox_notification_group.by_name.group_ids
  description = "List of user group IDs in this notification group"
}

output "notification_group_user_ids" {
  value       = data.netbox_notification_group.by_name.user_ids
  description = "List of user IDs in this notification group"
}

# Note: Notification groups do not support custom fields in NetBox API
output "notification_group_note" {
  value       = "Notification groups do not support custom fields"
  description = "Notification groups are administrative objects that manage notification recipients and cannot have custom fields assigned"
}
