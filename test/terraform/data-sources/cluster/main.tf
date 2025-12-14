# Cluster Data Source Test
# This test creates a cluster resource, then looks it up using the data source

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

# Dependency: Cluster Type
resource "netbox_cluster_type" "source_type" {
  name        = "DS Test Cluster Type"
  slug        = "ds-test-cluster-type"
  description = "Cluster type for data source testing"
}

# Create a cluster to look up
resource "netbox_cluster" "source" {
  name        = "Data Source Test Cluster"
  type        = netbox_cluster_type.source_type.slug
  status      = "active"
  description = "Cluster created for data source testing"
  comments    = "Test cluster for DS tests"

  depends_on = [netbox_cluster_type.source_type]
}

# Test 1: Look up by ID
data "netbox_cluster" "by_id" {
  id = netbox_cluster.source.id

  depends_on = [netbox_cluster.source]
}

# Test 2: Look up by name
data "netbox_cluster" "by_name" {
  name = netbox_cluster.source.name

  depends_on = [netbox_cluster.source]
}

# Outputs for verification
output "source_id" {
  value = netbox_cluster.source.id
}

output "by_id_name" {
  value = data.netbox_cluster.by_id.name
}

output "by_id_type" {
  value = data.netbox_cluster.by_id.type
}

output "by_name_status" {
  value = data.netbox_cluster.by_name.status
}

output "by_name_description" {
  value = data.netbox_cluster.by_name.description
}

# Verify all lookups return the same ID
output "ids_match" {
  value = data.netbox_cluster.by_id.id == data.netbox_cluster.by_name.id
}
