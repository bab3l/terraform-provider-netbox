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

# Dependencies
resource "netbox_site" "test" {
  name   = "Test Site for Rack Reservation DS"
  slug   = "test-site-rack-resv-ds"
  status = "active"
}

resource "netbox_rack" "test" {
  name     = "test-rack-reservation-ds"
  site     = netbox_site.test.id
  status   = "active"
  u_height = 42
}

data "netbox_user" "admin" {
  username = "admin"
}

resource "netbox_rack_reservation" "test" {
  rack        = netbox_rack.test.id
  units       = [6, 7]
  user        = data.netbox_user.admin.id
  description = "Rack reservation for data source test"
}

data "netbox_rack_reservation" "by_id" {
  id = netbox_rack_reservation.test.id
}
