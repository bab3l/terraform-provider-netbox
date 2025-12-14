# Circuit Type Data Source Test

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
resource "netbox_circuit_type" "test" {
  name        = "Test Circuit Type DS"
  slug        = "test-circuit-type-ds"
  description = "Test circuit type for data source"
}

# Test: Lookup circuit type by ID
data "netbox_circuit_type" "by_id" {
  id = netbox_circuit_type.test.id
}

# Test: Lookup circuit type by name
data "netbox_circuit_type" "by_name" {
  name = netbox_circuit_type.test.name

  depends_on = [netbox_circuit_type.test]
}
