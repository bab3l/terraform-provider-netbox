# Test to verify custom fields are preserved during updates
# This test creates a device with custom fields, then updates an unrelated field
# to verify that custom fields are not deleted during the update operation.

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

# Create device with custom fields
resource "netbox_device" "test" {
  name        = "test-device-cf"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id

  # Set custom fields
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

  description = "Initial description"
}

# Output to verify custom fields
output "device_id" {
  value = netbox_device.test.id
}

output "custom_fields" {
  value = netbox_device.test.custom_fields
}
