# Test: Config Context resource
# Tests the netbox_config_context resource CRUD operations

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

# Basic config context with minimal fields
resource "netbox_config_context" "basic" {
  name = "Test Config Basic ${var.test_id}"
  data = jsonencode({
    dns_servers = ["8.8.8.8", "8.8.4.4"]
  })
}

# Config context with all common fields
resource "netbox_config_context" "full" {
  name        = "Test Config Full ${var.test_id}"
  description = "Full config context for testing"
  weight      = 2000
  is_active   = true
  data = jsonencode({
    ntp_servers = ["time.google.com"]
    syslog = {
      server = "10.0.0.1"
      port   = 514
    }
    snmp = {
      community = "public"
      version   = "2c"
    }
  })
}

# Config context with weight and inactive status
resource "netbox_config_context" "inactive" {
  name        = "Test Config Inactive ${var.test_id}"
  description = "Inactive config context"
  weight      = 500
  is_active   = false
  data = jsonencode({
    test = "value"
  })
}

# Data source tests
data "netbox_config_context" "by_id" {
  id = netbox_config_context.basic.id
}

data "netbox_config_context" "by_name" {
  name = netbox_config_context.full.name
}

# Outputs for verification
output "basic_id_valid" {
  value = netbox_config_context.basic.id != ""
}

output "basic_name_valid" {
  value = netbox_config_context.basic.name == "Test Config Basic ${var.test_id}"
}

output "full_weight_valid" {
  value = netbox_config_context.full.weight == 2000
}

output "full_is_active_valid" {
  value = netbox_config_context.full.is_active == true
}

output "full_description_valid" {
  value = netbox_config_context.full.description == "Full config context for testing"
}

output "inactive_is_active_valid" {
  value = netbox_config_context.inactive.is_active == false
}

output "inactive_weight_valid" {
  value = netbox_config_context.inactive.weight == 500
}

output "datasource_by_id_matches" {
  value = data.netbox_config_context.by_id.name == netbox_config_context.basic.name
}

output "datasource_by_name_matches" {
  value = data.netbox_config_context.by_name.id == netbox_config_context.full.id
}
