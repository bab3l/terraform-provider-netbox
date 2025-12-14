# Look up a webhook by ID
data "netbox_webhook" "by_id" {
  id = "123"
}

# Look up a webhook by name
data "netbox_webhook" "by_name" {
  name = "slack-notification"
}

# Use webhook data in another configuration
output "webhook_url" {
  value = data.netbox_webhook.by_name.payload_url
}

output "webhook_method" {
  value = data.netbox_webhook.by_name.http_method
}

output "webhook_by_id" {
  value = data.netbox_webhook.by_id
}
