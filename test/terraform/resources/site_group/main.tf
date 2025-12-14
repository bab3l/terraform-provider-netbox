# Site Group Resource Test

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

# Test 1: Basic site group (root level)
resource "netbox_site_group" "root" {
  name        = "Test Root Site Group"
  slug        = "test-root-site-group"
  description = "A root-level site group for testing"
}

# Test 2: Child site group (with parent)
resource "netbox_site_group" "child" {
  name        = "Test Child Site Group"
  slug        = "test-child-site-group"
  parent      = netbox_site_group.root.id
  description = "A child site group under the root"
}

# Test 3: Grandchild site group (deeper hierarchy)
resource "netbox_site_group" "grandchild" {
  name        = "Test Grandchild Site Group"
  slug        = "test-grandchild-site-group"
  parent      = netbox_site_group.child.id
  description = "A grandchild site group for testing deep hierarchies"
}

# Test 4: Another root-level group (sibling)
resource "netbox_site_group" "sibling" {
  name        = "Test Sibling Site Group"
  slug        = "test-sibling-site-group"
  description = "A sibling root-level site group"
}
