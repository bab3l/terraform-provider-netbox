# Circuit Group Data Source Test

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

# Create a tenant to use for testing
resource "netbox_tenant" "test" {
  name = "test-circuit-group-ds-tenant"
  slug = "test-circuit-group-ds-tenant"
}

# Create a circuit group resource to test the data source
resource "netbox_circuit_group" "test" {
  name        = "test-circuit-group-ds"
  slug        = "test-circuit-group-ds"
  description = "Test circuit group for data source"
  tenant      = netbox_tenant.test.id
}

# Test 1: Look up by ID
data "netbox_circuit_group" "by_id" {
  id = netbox_circuit_group.test.id
}

# Test 2: Look up by slug
data "netbox_circuit_group" "by_slug" {
  slug = netbox_circuit_group.test.slug

  depends_on = [netbox_circuit_group.test]
}

# Test 3: Look up by name
data "netbox_circuit_group" "by_name" {
  name = netbox_circuit_group.test.name

  depends_on = [netbox_circuit_group.test]
}

# Outputs for verification
output "by_id_name" {
  value = data.netbox_circuit_group.by_id.name
}

output "by_id_slug" {
  value = data.netbox_circuit_group.by_id.slug
}

output "by_id_description" {
  value = data.netbox_circuit_group.by_id.description
}

output "by_id_tenant" {
  value = data.netbox_circuit_group.by_id.tenant
}

output "by_id_tenant_id" {
  value = data.netbox_circuit_group.by_id.tenant_id
}

output "by_slug_name" {
  value = data.netbox_circuit_group.by_slug.name
}

output "by_slug_id" {
  value = data.netbox_circuit_group.by_slug.id
}

output "by_name_id" {
  value = data.netbox_circuit_group.by_name.id
}

output "by_name_slug" {
  value = data.netbox_circuit_group.by_name.slug
}
