# Rack Role Resource Test

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

# Basic rack role with only required fields
resource "netbox_rack_role" "basic" {
  name = "Basic Test Rack Role"
  slug = "basic-test-rack-role"
}

# Complete rack role with all optional fields
resource "netbox_rack_role" "complete" {
  name        = "Complete Test Rack Role"
  slug        = "complete-test-rack-role"
  color       = "aa1409"
  description = "This is a complete test rack role with all fields."
}

# Rack role with specific color (e.g., for Production)
resource "netbox_rack_role" "production" {
  name        = "Production Rack Role"
  slug        = "production-rack-role"
  color       = "2ecc71"
  description = "Racks designated for production workloads."
}

# Rack role for testing environments
resource "netbox_rack_role" "testing" {
  name        = "Testing Rack Role"
  slug        = "testing-rack-role"
  color       = "f39c12"
  description = "Racks designated for testing and staging."
}

# Rack role for storage
resource "netbox_rack_role" "storage" {
  name        = "Storage Rack Role"
  slug        = "storage-rack-role"
  color       = "3498db"
  description = "Racks designated for storage systems."
}

# Output values for verification
output "basic_id" {
  value = netbox_rack_role.basic.id
}

output "basic_name" {
  value = netbox_rack_role.basic.name
}

output "basic_slug" {
  value = netbox_rack_role.basic.slug
}

output "complete_id" {
  value = netbox_rack_role.complete.id
}

output "complete_name" {
  value = netbox_rack_role.complete.name
}

output "complete_slug" {
  value = netbox_rack_role.complete.slug
}

output "complete_color" {
  value = netbox_rack_role.complete.color
}

output "complete_description" {
  value = netbox_rack_role.complete.description
}

output "production_id" {
  value = netbox_rack_role.production.id
}

output "production_color" {
  value = netbox_rack_role.production.color
}

output "testing_id" {
  value = netbox_rack_role.testing.id
}

output "testing_color" {
  value = netbox_rack_role.testing.color
}

output "storage_id" {
  value = netbox_rack_role.storage.id
}

output "storage_color" {
  value = netbox_rack_role.storage.color
}
