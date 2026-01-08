# Phase 2: Update description WITHOUT mentioning custom_fields
# This simulates a user updating an unrelated field
# Custom fields should be preserved, not deleted

terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  server_url = "http://localhost:8000"
  api_token  = "0123456789abcdef0123456789abcdef01234567"
}

# Create custom fields for device
resource "netbox_custom_field" "test_text" {
  name = "test_text_field"
  type = "text"
  content_types = ["dcim.device"]
  description = "Test text custom field"
}

resource "netbox_custom_field" "test_integer" {
  name = "test_integer_field"
  type = "integer"
  content_types = ["dcim.device"]
  description = "Test integer custom field"
}

# Create dependencies
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_type" "test" {
  model        = "Test Model"
  slug         = "test-model"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = "Test Role"
  slug  = "test-role"
  color = "ff0000"
}

resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

# Update device - change description but keep custom fields in config
resource "netbox_device" "test" {
  name        = "test-device-cf"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id

  # Keep custom fields the same
  custom_fields = [
    {
      name  = netbox_custom_field.test_text.name
      type  = "text"
      value = "initial value"
    },
    {
      name  = netbox_custom_field.test_integer.name
      type  = "integer"
      value = "42"
    }
  ]

  # Update this field
  description = "Updated description - custom fields should remain!"
}

# Output to verify custom fields are still present
output "device_id" {
  value = netbox_device.test.id
}

output "custom_fields" {
  value = netbox_device.test.custom_fields
}

output "description" {
  value = netbox_device.test.description
}
