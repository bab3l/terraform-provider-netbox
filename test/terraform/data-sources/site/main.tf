# Site Data Source Test
# This test creates a site resource, then looks it up using the data source

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

# First, create a site to look up
resource "netbox_site" "source" {
  name        = "Data Source Test Site"
  slug        = "data-source-test-site"
  status      = "active"
  facility    = "FAC-DS-01"
  description = "Site created for data source testing"
}

# Test 1: Look up by ID
data "netbox_site" "by_id" {
  id = netbox_site.source.id

  depends_on = [netbox_site.source]
}

# Test 2: Look up by name
data "netbox_site" "by_name" {
  name = netbox_site.source.name

  depends_on = [netbox_site.source]
}

# Test 3: Look up by slug
data "netbox_site" "by_slug" {
  slug = netbox_site.source.slug

  depends_on = [netbox_site.source]
}
