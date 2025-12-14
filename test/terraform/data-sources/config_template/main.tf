# Config Template Data Source Test

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

# Dependencies
resource "netbox_config_template" "test" {
  name            = "Test Config Template DS"
  template_code   = "hostname {{ hostname }}"
  description     = "Test config template for data source"
}

# Test: Lookup config template by ID
data "netbox_config_template" "by_id" {
  id = netbox_config_template.test.id
}

# Test: Lookup config template by name
data "netbox_config_template" "by_name" {
  name = netbox_config_template.test.name

  depends_on = [netbox_config_template.test]
}
