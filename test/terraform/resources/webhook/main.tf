# Test: Webhook resource
# Tests the netbox_webhook resource CRUD operations

terraform {
  required_version = ">= 1.0"
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

# Basic webhook
resource "netbox_webhook" "basic" {
  name        = "Test Webhook Basic ${var.test_id}"
  payload_url = "https://example.com/webhook/basic"
}

# Webhook with all fields
resource "netbox_webhook" "full" {
  name               = "Test Webhook Full ${var.test_id}"
  payload_url        = "https://example.com/webhook/full"
  http_method        = "PUT"
  http_content_type  = "application/xml"
  description        = "Test webhook with all fields"
  additional_headers = "X-Custom-Header: test-value"
  ssl_verification   = false
}

# Data source tests
data "netbox_webhook" "by_id" {
  id = netbox_webhook.basic.id
}

data "netbox_webhook" "by_name" {
  name = netbox_webhook.full.name
}

# Outputs for verification
output "basic_webhook_id" {
  value = netbox_webhook.basic.id
}

output "full_webhook_id" {
  value = netbox_webhook.full.id
}

output "datasource_by_id_name" {
  value = data.netbox_webhook.by_id.name
}

output "datasource_by_name_url" {
  value = data.netbox_webhook.by_name.payload_url
}

output "datasource_by_name_method" {
  value = data.netbox_webhook.by_name.http_method
}
