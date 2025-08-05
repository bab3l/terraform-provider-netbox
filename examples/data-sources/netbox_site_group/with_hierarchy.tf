# Site Group Data Source with Hierarchical Organization

terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  server_url = "https://netbox.example.com"
  api_token  = "your-api-token-here"
}

# Create a parent site group
resource "netbox_site_group" "regional_group" {
  name        = "North America"
  slug        = "north-america"
  description = "Sites located in North America"

  tags = [
    {
      name = "region"
      slug = "region"
    }
  ]

  custom_fields = [
    {
      name  = "region_code"
      type  = "text"
      value = "NA"
    }
  ]
}

# Create a child site group
resource "netbox_site_group" "country_group" {
  name        = "United States"
  slug        = "united-states"
  parent      = netbox_site_group.regional_group.id
  description = "Sites located in the United States"

  tags = [
    {
      name = "country"
      slug = "country"
    }
  ]

  custom_fields = [
    {
      name  = "country_code"
      type  = "text"
      value = "US"
    }
  ]
}

# Use data sources to read the created site groups
data "netbox_site_group" "regional_data" {
  id = netbox_site_group.regional_group.id
}

data "netbox_site_group" "country_data" {
  slug = netbox_site_group.country_group.slug
}

# Read an existing site group to demonstrate lookup
data "netbox_site_group" "existing_group" {
  slug = "production-sites"
}

# Output hierarchical information
output "site_group_hierarchy" {
  description = "Information about the site group hierarchy"
  value = {
    parent_group = {
      id            = data.netbox_site_group.regional_data.id
      name          = data.netbox_site_group.regional_data.name
      slug          = data.netbox_site_group.regional_data.slug
      parent        = data.netbox_site_group.regional_data.parent
      description   = data.netbox_site_group.regional_data.description
      tags          = data.netbox_site_group.regional_data.tags
      custom_fields = data.netbox_site_group.regional_data.custom_fields
    }

    child_group = {
      id            = data.netbox_site_group.country_data.id
      name          = data.netbox_site_group.country_data.name
      slug          = data.netbox_site_group.country_data.slug
      parent        = data.netbox_site_group.country_data.parent
      description   = data.netbox_site_group.country_data.description
      tags          = data.netbox_site_group.country_data.tags
      custom_fields = data.netbox_site_group.country_data.custom_fields
    }

    existing_group = {
      id          = data.netbox_site_group.existing_group.id
      name        = data.netbox_site_group.existing_group.name
      slug        = data.netbox_site_group.existing_group.slug
      parent      = data.netbox_site_group.existing_group.parent
      description = data.netbox_site_group.existing_group.description
    }
  }
}

# Example of using site group data for conditional logic
locals {
  is_top_level_group = data.netbox_site_group.regional_data.parent == null
  has_child_group    = data.netbox_site_group.country_data.parent != null
  parent_matches     = data.netbox_site_group.country_data.parent == data.netbox_site_group.regional_data.name
}

output "group_relationships" {
  value = {
    regional_is_top_level       = local.is_top_level_group
    country_has_parent          = local.has_child_group
    parent_relationship_correct = local.parent_matches
    message                     = local.parent_matches ? "Hierarchy is correctly configured" : "Check parent-child relationships"
  }
}
