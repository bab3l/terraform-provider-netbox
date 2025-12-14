# Route Target Resource Test

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

# Test 1: Basic route target creation
resource "netbox_route_target" "basic" {
  name = "65000:100"
}

# Test 2: Route target with description and comments
resource "netbox_route_target" "complete" {
  name        = "65000:200"
  description = "Test route target for VRF export"
  comments    = "This is a test route target"
}

# Test 3: Route target with tenant
resource "netbox_tenant" "rt_test" {
  name = "RT Test Tenant"
  slug = "rt-test-tenant"
}

resource "netbox_route_target" "with_tenant" {
  name        = "65000:300"
  tenant      = netbox_tenant.rt_test.id
  description = "Route target owned by tenant"
}
