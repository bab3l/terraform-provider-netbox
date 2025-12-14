# rack_reservation resource test
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
resource "netbox_site" "test" {
  name   = "Test Site for Rack Reservation"
  slug   = "test-site-rack-reservation"
  status = "active"
}

resource "netbox_rack" "test" {
  name     = "test-rack-reservation"
  site     = netbox_site.test.id
  status   = "active"
  u_height = 42
}

data "netbox_user" "admin" {
  username = "admin"
}

# Test 1: Basic rack reservation with units
resource "netbox_rack_reservation" "basic" {
  rack        = netbox_rack.test.id
  units       = [1, 2]
  user        = data.netbox_user.admin.id
  description = "Basic rack reservation"
}

# Test 2: Rack reservation with description
resource "netbox_rack_reservation" "complete" {
  rack        = netbox_rack.test.id
  units       = [3, 4, 5]
  user        = data.netbox_user.admin.id
  description = "Test rack reservation with full details"
}
