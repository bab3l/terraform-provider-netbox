# Circuit Group Assignment Data Source Test

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

# Create prerequisite resources

# Provider
resource "netbox_provider" "test" {
  name = "test-circuit-group-assignment-ds-provider"
  slug = "test-circuit-group-assignment-ds-provider"
}

# Circuit type
resource "netbox_circuit_type" "test" {
  name = "test-circuit-group-assignment-ds-type"
  slug = "test-circuit-group-assignment-ds-type"
}

# Circuit group
resource "netbox_circuit_group" "test" {
  name = "test-circuit-group-assignment-ds-group"
  slug = "test-circuit-group-assignment-ds-group"
}

# Circuit
resource "netbox_circuit" "test" {
  cid              = "TEST-CGA-DS-001"
  circuit_provider = netbox_provider.test.slug
  type             = netbox_circuit_type.test.slug
}

# Create a circuit group assignment resource to test the data source
resource "netbox_circuit_group_assignment" "test" {
  group_id   = netbox_circuit_group.test.id
  circuit_id = netbox_circuit.test.id
  priority   = "primary"
}

# Test 1: Look up by ID
data "netbox_circuit_group_assignment" "by_id" {
  id = netbox_circuit_group_assignment.test.id
}

# Outputs for verification
output "by_id_id" {
  value = data.netbox_circuit_group_assignment.by_id.id
}

output "by_id_group_id" {
  value = data.netbox_circuit_group_assignment.by_id.group_id
}

output "by_id_group_name" {
  value = data.netbox_circuit_group_assignment.by_id.group_name
}

output "by_id_circuit_id" {
  value = data.netbox_circuit_group_assignment.by_id.circuit_id
}

output "by_id_circuit_cid" {
  value = data.netbox_circuit_group_assignment.by_id.circuit_cid
}

output "by_id_priority" {
  value = data.netbox_circuit_group_assignment.by_id.priority
}

output "by_id_priority_name" {
  value = data.netbox_circuit_group_assignment.by_id.priority_name
}
