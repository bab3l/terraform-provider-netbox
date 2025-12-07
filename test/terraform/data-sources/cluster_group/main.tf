# Cluster Group Data Source Test

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

# Dependencies
resource "netbox_cluster_group" "test" {
  name        = "Test Cluster Group DS"
  slug        = "test-cluster-group-ds"
  description = "Test cluster group for data source"
}

# Test: Lookup cluster group by ID
data "netbox_cluster_group" "by_id" {
  id = netbox_cluster_group.test.id
}

# Test: Lookup cluster group by name
data "netbox_cluster_group" "by_name" {
  name = netbox_cluster_group.test.name

  depends_on = [netbox_cluster_group.test]
}
