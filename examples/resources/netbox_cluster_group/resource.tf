# Example: Basic cluster group
resource "netbox_cluster_group" "basic" {
  name = "Production Clusters"
  slug = "production-clusters"
}

# Example: Cluster group with description
resource "netbox_cluster_group" "development" {
  name        = "Development Clusters"
  slug        = "development-clusters"
  description = "Clusters used for development environments"
}

# Example: Cluster group for staging
resource "netbox_cluster_group" "staging" {
  name        = "Staging Clusters"
  slug        = "staging-clusters"
  description = "Clusters used for staging and QA environments"
}

# Example: Cluster group by datacenter
resource "netbox_cluster_group" "datacenter_us" {
  name        = "US Datacenters"
  slug        = "us-datacenters"
  description = "All clusters in US datacenters"
}

resource "netbox_cluster_group" "datacenter_eu" {
  name        = "EU Datacenters"
  slug        = "eu-datacenters"
  description = "All clusters in EU datacenters"
}

# Example: Cluster group with tags
resource "netbox_cluster_group" "with_tags" {
  name        = "Critical Infrastructure"
  slug        = "critical-infrastructure"
  description = "Mission-critical virtualization clusters"
  tags {
    name = "critical"
    slug = "critical"
  }
}
