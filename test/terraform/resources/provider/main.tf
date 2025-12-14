# Provider Integration Test
# Tests the netbox_provider resource with basic and complete configurations
# Note: This is a circuit provider (ISP/carrier), not a Terraform provider

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

# Basic Provider with only required fields
resource "netbox_provider" "basic" {
  name = "Basic Test Provider"
  slug = "basic-test-provider"
}

# Complete Provider with all optional fields
resource "netbox_provider" "complete" {
  name        = "Complete Test Provider"
  slug        = "complete-test-provider"
  description = "Complete provider for integration testing"
  comments    = "Created by terraform integration test"
}

# Provider for ASN testing (if ASN resource is available)
resource "netbox_provider" "with_details" {
  name        = "ISP Provider"
  slug        = "isp-provider"
  description = "Internet Service Provider for testing"
  comments    = "ISP details and configuration"
}
