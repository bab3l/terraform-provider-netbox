# Look up an existing notification group by ID
data "netbox_notification_group" "ops_team" {
  id = "1"
}

# Use the notification group data
output "notification_group_name" {
  value = data.netbox_notification_group.ops_team.name
}

output "notification_group_description" {
  value = data.netbox_notification_group.ops_team.description
}
