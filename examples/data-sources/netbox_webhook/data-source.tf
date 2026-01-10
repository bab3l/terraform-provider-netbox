# Look up a webhook by ID
data "netbox_webhook" "by_id" {
  id = "123"
}

# Look up a webhook by name
data "netbox_webhook" "by_name" {
  name = "slack-notification"
}

# Use webhook data in another configuration
output "webhook_id" {
  value = data.netbox_webhook.by_id.id
}

output "webhook_name" {
  value = data.netbox_webhook.by_name.name
}

output "webhook_url" {
  value = data.netbox_webhook.by_name.payload_url
}

output "webhook_method" {
  value = data.netbox_webhook.by_name.http_method
}

output "webhook_ssl_verification" {
  value = data.netbox_webhook.by_id.ssl_verification
}

output "webhook_secret" {
  value     = data.netbox_webhook.by_id.secret
  sensitive = true
}

# Note: Webhooks do not support custom fields in NetBox API
output "webhook_event_types" {
  value       = data.netbox_webhook.by_id.events
  description = "Events that trigger this webhook"
}

output "webhook_content_type" {
  value       = data.netbox_webhook.by_id.content_type
  description = "Content type for webhook payload"
}
