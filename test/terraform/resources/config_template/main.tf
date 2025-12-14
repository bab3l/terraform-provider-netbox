# Config Template Resource Test

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

# Test 1: Basic config template creation
resource "netbox_config_template" "basic" {
  name            = "Test Config Template Basic"
  template_code   = "hostname {{ device.name }}"
}

# Test 2: Config template with all optional fields
resource "netbox_config_template" "complete" {
  name            = "Test Config Template Complete"
  description     = "A complete test config template"
  template_code   = "hostname {{ device.name }}\ninterface Loopback0\n  ip address {{ device.primary_ip }}"
}
