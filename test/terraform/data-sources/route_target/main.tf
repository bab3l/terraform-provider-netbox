# Route Target Data Source Test

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

# Create a route target to look up
resource "netbox_route_target" "test" {
  name        = "65001:100"
  description = "Test route target for data source"
}

# Look up by ID
data "netbox_route_target" "by_id" {
  id = netbox_route_target.test.id
}

# Look up by name
data "netbox_route_target" "by_name" {
  name = netbox_route_target.test.name
}
