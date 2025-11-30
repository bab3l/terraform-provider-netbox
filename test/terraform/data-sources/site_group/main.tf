# Site Group Data Source Test

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

# First, create site groups to look up
resource "netbox_site_group" "source_parent" {
  name        = "DS Test Parent Group"
  slug        = "ds-test-parent-group"
  description = "Parent site group for data source testing"
}

resource "netbox_site_group" "source_child" {
  name        = "DS Test Child Group"
  slug        = "ds-test-child-group"
  parent      = netbox_site_group.source_parent.id
  description = "Child site group for data source testing"
}

# Test 1: Look up parent by ID
data "netbox_site_group" "by_id" {
  id = netbox_site_group.source_parent.id

  depends_on = [netbox_site_group.source_parent]
}

# Test 2: Look up by name
data "netbox_site_group" "by_name" {
  name = netbox_site_group.source_parent.name

  depends_on = [netbox_site_group.source_parent]
}

# Test 3: Look up by slug
data "netbox_site_group" "by_slug" {
  slug = netbox_site_group.source_parent.slug

  depends_on = [netbox_site_group.source_parent]
}

# Test 4: Look up child (has parent)
data "netbox_site_group" "child_by_id" {
  id = netbox_site_group.source_child.id

  depends_on = [netbox_site_group.source_child]
}
