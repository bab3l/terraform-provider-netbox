# Test: Config Context data source
# This tests looking up config contexts by various identifiers

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

# Create a config context to look up
resource "netbox_config_context" "test" {
  name        = "Test Config DS ${var.test_id}"
  description = "Config context for data source testing"
  weight      = 1500
  is_active   = true
  data = jsonencode({
    test_key = "test_value"
    nested = {
      key1 = "value1"
      key2 = "value2"
    }
  })
}

# Look up by ID
data "netbox_config_context" "by_id" {
  id = netbox_config_context.test.id
}

# Look up by name
data "netbox_config_context" "by_name" {
  name = netbox_config_context.test.name
}

# Verification outputs
output "id_matches" {
  value = data.netbox_config_context.by_id.id == netbox_config_context.test.id
}

output "name_matches" {
  value = data.netbox_config_context.by_name.name == netbox_config_context.test.name
}

output "weight_matches" {
  value = data.netbox_config_context.by_id.weight == netbox_config_context.test.weight
}

output "is_active_matches" {
  value = data.netbox_config_context.by_id.is_active == netbox_config_context.test.is_active
}

output "description_matches" {
  value = data.netbox_config_context.by_name.description == netbox_config_context.test.description
}
