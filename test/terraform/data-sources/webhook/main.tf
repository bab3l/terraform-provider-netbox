# Test: Webhook data source
# This tests looking up webhooks by various identifiers

terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  # Uses NETBOX_SERVER_URL and NETBOX_API_TOKEN environment variables
}

variable "test_id" {
  description = "Unique identifier for this test run"
  type        = string
  default     = "test"
}

# Create a webhook to look up
resource "netbox_webhook" "test" {
  name        = "Test Webhook DS ${var.test_id}"
  payload_url = "https://example.com/webhook/ds-test"
  http_method = "POST"
  description = "Test webhook for data source testing"
}

# Look up by ID
data "netbox_webhook" "by_id" {
  id = netbox_webhook.test.id
}

# Look up by name
data "netbox_webhook" "by_name" {
  name = netbox_webhook.test.name
}

# Verification outputs
output "lookup_by_id_name" {
  value = data.netbox_webhook.by_id.name
}

output "lookup_by_name_url" {
  value = data.netbox_webhook.by_name.payload_url
}

output "id_matches" {
  value = data.netbox_webhook.by_id.id == netbox_webhook.test.id
}

output "name_matches" {
  value = data.netbox_webhook.by_name.name == netbox_webhook.test.name
}
