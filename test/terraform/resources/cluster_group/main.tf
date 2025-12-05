# Cluster Group Resource Integration Tests
# Tests the netbox_cluster_group resource CRUD operations

terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {}

# Test 1: Basic cluster group with minimal required fields
resource "netbox_cluster_group" "test_basic" {
  name = "Test Cluster Group Basic"
  slug = "test-cluster-group-basic"
}

# Test 2: Cluster group with all optional fields
resource "netbox_cluster_group" "test_full" {
  name        = "Test Cluster Group Full"
  slug        = "test-cluster-group-full"
  description = "A fully configured cluster group for integration testing"
}

# Test 3: Cluster group for production workloads
resource "netbox_cluster_group" "test_production" {
  name        = "Production Clusters"
  slug        = "production-clusters"
  description = "Cluster group for production workloads"
}

# Test 4: Cluster group for development
resource "netbox_cluster_group" "test_development" {
  name        = "Development Clusters"
  slug        = "development-clusters"
  description = "Cluster group for development and testing"
}

# Data source test - lookup by ID
data "netbox_cluster_group" "by_id" {
  id = netbox_cluster_group.test_basic.id
}

# Data source test - lookup by name
data "netbox_cluster_group" "by_name" {
  name = netbox_cluster_group.test_full.name
  
  depends_on = [netbox_cluster_group.test_full]
}

# Data source test - lookup by slug
data "netbox_cluster_group" "by_slug" {
  slug = netbox_cluster_group.test_production.slug
  
  depends_on = [netbox_cluster_group.test_production]
}

# Outputs for verification
output "basic_group_id_valid" {
  value = can(tonumber(netbox_cluster_group.test_basic.id))
}

output "basic_group_name_valid" {
  value = netbox_cluster_group.test_basic.name == "Test Cluster Group Basic"
}

output "full_group_description_valid" {
  value = netbox_cluster_group.test_full.description == "A fully configured cluster group for integration testing"
}

output "production_group_name_valid" {
  value = netbox_cluster_group.test_production.name == "Production Clusters"
}

output "development_group_slug_valid" {
  value = netbox_cluster_group.test_development.slug == "development-clusters"
}

output "data_by_id_matches" {
  value = data.netbox_cluster_group.by_id.name == netbox_cluster_group.test_basic.name
}

output "data_by_name_matches" {
  value = data.netbox_cluster_group.by_name.slug == netbox_cluster_group.test_full.slug
}

output "data_by_slug_matches" {
  value = data.netbox_cluster_group.by_slug.description == netbox_cluster_group.test_production.description
}
