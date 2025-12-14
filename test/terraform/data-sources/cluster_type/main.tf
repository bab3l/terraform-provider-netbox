# Cluster Type Data Source Test
# This test creates a cluster type resource, then looks it up using the data source

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

# First, create a cluster type to look up
resource "netbox_cluster_type" "source" {
  name        = "Data Source Test Cluster Type"
  slug        = "data-source-test-cluster-type"
  description = "Cluster type created for data source testing"
}

# Test 1: Look up by ID
data "netbox_cluster_type" "by_id" {
  id = netbox_cluster_type.source.id

  depends_on = [netbox_cluster_type.source]
}

# Test 2: Look up by name
data "netbox_cluster_type" "by_name" {
  name = netbox_cluster_type.source.name

  depends_on = [netbox_cluster_type.source]
}

# Test 3: Look up by slug
data "netbox_cluster_type" "by_slug" {
  slug = netbox_cluster_type.source.slug

  depends_on = [netbox_cluster_type.source]
}

# Outputs for verification
output "source_id" {
  value = netbox_cluster_type.source.id
}

output "by_id_name" {
  value = data.netbox_cluster_type.by_id.name
}

output "by_name_slug" {
  value = data.netbox_cluster_type.by_name.slug
}

output "by_slug_description" {
  value = data.netbox_cluster_type.by_slug.description
}

# Verify all lookups return the same ID
output "all_ids_match" {
  value = (
    data.netbox_cluster_type.by_id.id == data.netbox_cluster_type.by_name.id &&
    data.netbox_cluster_type.by_name.id == data.netbox_cluster_type.by_slug.id
  )
}
