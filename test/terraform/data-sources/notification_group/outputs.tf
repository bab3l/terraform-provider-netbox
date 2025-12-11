output "notification_group_id" {
  value = data.netbox_notification_group.test.id
}

output "notification_group_name" {
  value = data.netbox_notification_group.test.name
}

output "notification_group_description" {
  value = data.netbox_notification_group.test.description
}
