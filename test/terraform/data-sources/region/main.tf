# Region Data Source Test

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

# First, create a region to look up
resource "netbox_region" "test_region" {
  name        = "Data Source Test Region"
  slug        = "data-source-test-region"
  description = "A region for testing data source lookups"
}

# Create a child region for hierarchy testing
resource "netbox_region" "child_region" {
  name        = "Data Source Child Region"
  slug        = "data-source-child-region"
  parent      = netbox_region.test_region.id
  description = "A child region for testing data source lookups"
}

# Data source: look up region by ID
data "netbox_region" "by_id" {
  id = netbox_region.test_region.id

  depends_on = [netbox_region.test_region]
}

# Data source: look up region by name
data "netbox_region" "by_name" {
  name = netbox_region.test_region.name

  depends_on = [netbox_region.test_region]
}

# Data source: look up region by slug
data "netbox_region" "by_slug" {
  slug = netbox_region.test_region.slug

  depends_on = [netbox_region.test_region]
}

# Data source: look up child region to verify parent relationship
data "netbox_region" "child_by_id" {
  id = netbox_region.child_region.id

  depends_on = [netbox_region.child_region]
}
