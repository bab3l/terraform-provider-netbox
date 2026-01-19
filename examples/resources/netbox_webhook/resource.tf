# Example: Basic webhook
resource "netbox_webhook" "basic" {
  name        = "slack-notification"
  payload_url = "https://hooks.slack.com/services/xxx/yyy/zzz"
}

# Example: Webhook with custom HTTP method and content type
resource "netbox_webhook" "custom" {
  name              = "custom-webhook"
  payload_url       = "https://api.example.com/webhook"
  http_method       = "PUT"
  http_content_type = "application/xml"
  description       = "Custom webhook with XML payload"
}

# Example: Webhook with additional headers and secret
resource "netbox_webhook" "secure" {
  name               = "secure-webhook"
  payload_url        = "https://secure.example.com/webhook"
  secret             = "my-secret-key"
  additional_headers = "X-API-Key: my-api-key\nX-Source: netbox"
  ssl_verification   = true
}

# Example: Webhook with custom body template
resource "netbox_webhook" "templated" {
  name        = "templated-webhook"
  payload_url = "https://api.example.com/events"
  body_template = jsonencode({
    event     = "{{ event }}"
    timestamp = "{{ timestamp }}"
    model     = "{{ model }}"
    data      = "{{ data | tojson }}"
  })
}

# Example: Webhook with SSL disabled (for testing only!)
resource "netbox_webhook" "insecure" {
  name             = "test-webhook"
  payload_url      = "http://localhost:8080/webhook"
  ssl_verification = false
  description      = "Test webhook with SSL disabled"
}

# Optional: seed owned custom fields during import
import {
  to = netbox_webhook.insecure
  id = "123"

  identity = {
    custom_fields = [
      "owner_team:text",
      "delivery_target:text",
    ]
  }
}
