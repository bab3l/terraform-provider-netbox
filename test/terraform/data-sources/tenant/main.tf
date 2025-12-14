# Tenant Data Source Test

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

# Create tenant group for testing
resource "netbox_tenant_group" "source_group" {
  name        = "DS Tenant Test Group"
  slug        = "ds-tenant-test-group"
  description = "Group for tenant data source testing"
}

# Create tenant to look up
resource "netbox_tenant" "source" {
  name        = "DS Test Tenant"
  slug        = "ds-test-tenant"
  group       = netbox_tenant_group.source_group.id
  description = "Tenant created for data source testing"
}

# Test 1: Look up by ID
data "netbox_tenant" "by_id" {
  id = netbox_tenant.source.id

  depends_on = [netbox_tenant.source]
}

# Test 2: Look up by name
data "netbox_tenant" "by_name" {
  name = netbox_tenant.source.name

  depends_on = [netbox_tenant.source]
}

# Test 3: Look up by slug
data "netbox_tenant" "by_slug" {
  slug = netbox_tenant.source.slug

  depends_on = [netbox_tenant.source]
}
