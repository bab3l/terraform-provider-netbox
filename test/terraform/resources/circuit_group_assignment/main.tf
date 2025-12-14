# Circuit Group Assignment Resource Test

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

# Create prerequisite resources

# Provider
resource "netbox_provider" "test" {
  name = "test-circuit-group-assignment-provider"
  slug = "test-circuit-group-assignment-provider"
}

# Circuit type
resource "netbox_circuit_type" "test" {
  name = "test-circuit-group-assignment-type"
  slug = "test-circuit-group-assignment-type"
}

# Circuit group
resource "netbox_circuit_group" "test" {
  name = "test-circuit-group-assignment-group"
  slug = "test-circuit-group-assignment-group"
}

# Circuits
resource "netbox_circuit" "primary" {
  cid              = "TEST-CGA-001"
  circuit_provider = netbox_provider.test.slug
  type             = netbox_circuit_type.test.slug
}

resource "netbox_circuit" "secondary" {
  cid              = "TEST-CGA-002"
  circuit_provider = netbox_provider.test.slug
  type             = netbox_circuit_type.test.slug
}

resource "netbox_circuit" "tertiary" {
  cid              = "TEST-CGA-003"
  circuit_provider = netbox_provider.test.slug
  type             = netbox_circuit_type.test.slug
}

# Test 1: Basic circuit group assignment
resource "netbox_circuit_group_assignment" "basic" {
  group_id   = netbox_circuit_group.test.id
  circuit_id = netbox_circuit.primary.id
}

# Test 2: Circuit group assignment with priority
resource "netbox_circuit_group_assignment" "primary" {
  group_id   = netbox_circuit_group.test.id
  circuit_id = netbox_circuit.secondary.id
  priority   = "primary"
}

# Test 3: Circuit group assignment with different priority
resource "netbox_circuit_group_assignment" "tertiary" {
  group_id   = netbox_circuit_group.test.id
  circuit_id = netbox_circuit.tertiary.id
  priority   = "tertiary"
}

# Test 4: Output values for verification
output "basic_assignment_id" {
  value = netbox_circuit_group_assignment.basic.id
}

output "basic_assignment_group_id" {
  value = netbox_circuit_group_assignment.basic.group_id
}

output "basic_assignment_circuit_id" {
  value = netbox_circuit_group_assignment.basic.circuit_id
}

output "primary_assignment_id" {
  value = netbox_circuit_group_assignment.primary.id
}

output "primary_assignment_priority" {
  value = netbox_circuit_group_assignment.primary.priority
}

output "tertiary_assignment_priority" {
  value = netbox_circuit_group_assignment.tertiary.priority
}
