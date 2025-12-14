# Custom Field Resource Test

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

# Test 1: Basic custom field creation
resource "netbox_custom_field" "basic" {
  name          = "tf_test_basic"
  type          = "text"
  object_types  = ["dcim.site"]
}

# Test 2: Custom field with all optional fields
resource "netbox_custom_field" "complete" {
  name          = "tf_test_complete"
  type          = "integer"
  object_types  = ["dcim.site", "dcim.device"]
  label         = "Test Complete Field"
  description   = "A complete test custom field"
  required      = false
  filter_logic  = "loose"
  ui_visible    = "always"
  ui_editable   = "yes"
  weight        = 100
}
