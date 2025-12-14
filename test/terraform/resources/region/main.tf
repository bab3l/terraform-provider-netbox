# Region Resource Test

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

# Test 1: Basic region creation
resource "netbox_region" "basic" {
  name        = "Test Region Basic"
  slug        = "test-region-basic"
  description = "A basic test region created by Terraform integration tests"
}

# Test 2: Region with all optional fields
resource "netbox_region" "complete" {
  name        = "Test Region Complete"
  slug        = "test-region-complete"
  description = "A complete test region with all available fields populated"
}

# Test 3: Nested regions (parent/child hierarchy)
resource "netbox_region" "parent" {
  name        = "Parent Region"
  slug        = "parent-region"
  description = "A parent region for testing hierarchies"
}

resource "netbox_region" "child" {
  name        = "Child Region"
  slug        = "child-region"
  parent      = netbox_region.parent.id
  description = "A child region nested under the parent"

  depends_on = [netbox_region.parent]
}

# Test 4: Deep nesting (grandchild)
resource "netbox_region" "grandchild" {
  name        = "Grandchild Region"
  slug        = "grandchild-region"
  parent      = netbox_region.child.id
  description = "A grandchild region for deep hierarchy testing"

  depends_on = [netbox_region.child]
}
