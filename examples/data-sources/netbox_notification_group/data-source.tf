# Look up an existing notification group by ID
data "netbox_notification_group" "ops_team" {
  id = "1"
}

# Use the notification group data
output "notification_group_name" {
  value       = data.netbox_notification_group.ops_team.name
  description = "Name of the notification group"
}

output "notification_group_description" {
  value       = data.netbox_notification_group.ops_team.description
  description = "Description of the notification group"
}

output "group_ids" {
  value       = data.netbox_notification_group.ops_team.group_ids
  description = "User group IDs in this notification group"
}

output "user_ids" {
  value       = data.netbox_notification_group.ops_team.user_ids
  description = "User IDs in this notification group"
}
