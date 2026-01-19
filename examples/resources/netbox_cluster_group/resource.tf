# Example: Basic cluster group
resource "netbox_cluster_group" "basic" {
  name = "Production Clusters"
  slug = "production-clusters"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "environment_type"
      value = "production"
    },
    {
      name  = "sla_tier"
      value = "tier-1"
    }
  ]

  tags = [
    "cluster-group",
    "production"
  ]
}

# Example: Cluster group with description
resource "netbox_cluster_group" "development" {
  name        = "Development Clusters"
  slug        = "development-clusters"
  description = "Clusters used for development environments"

  # Partial custom fields management
  custom_fields = [
    {
      name  = "environment_type"
      value = "development"
    },
    {
      name  = "auto_shutdown_enabled"
      value = "true"
    }
  ]

  tags = [
    "cluster-group",
    "development"
  ]
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
  tags        = ["critical"]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_cluster_group.basic
  id = "123"

  identity = {
    custom_fields = [
      "environment_type:text",
      "sla_tier:text",
    ]
  }
}

# Optional: seed owned custom fields during import
import {
  to = netbox_cluster_group.development
  id = "124"

  identity = {
    custom_fields = [
      "environment_type:text",
      "auto_shutdown_enabled:boolean",
    ]
  }
}
