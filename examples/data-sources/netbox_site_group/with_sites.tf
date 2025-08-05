# Site Group and Site Integration Example

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

# Create a site group for production sites
resource "netbox_site_group" "production_sites" {
  name        = "Production Sites"
  slug        = "production-sites"
  description = "All production data centers and facilities"

  tags = [
    {
      name = "production"
      slug = "production"
    },
    {
      name = "critical"
      slug = "critical"
    }
  ]
}

# Create a site that belongs to the site group
resource "netbox_site" "primary_datacenter" {
  name        = "Primary Datacenter"
  slug        = "primary-dc"
  status      = "active"
  group       = netbox_site_group.production_sites.id
  description = "Main production facility"
  facility    = "Building A"

  tags = [
    {
      name = "production"
      slug = "production"
    },
    {
      name = "primary"
      slug = "primary"
    }
  ]
}

# Use data sources to read the site group and site
data "netbox_site_group" "production_group_data" {
  id = netbox_site_group.production_sites.id
}

data "netbox_site" "primary_dc_data" {
  id = netbox_site.primary_datacenter.id
}

# Read existing site groups to find sites
data "netbox_site_group" "existing_dr_group" {
  slug = "disaster-recovery"
}

# Read an existing site that might belong to a group
data "netbox_site" "existing_site" {
  slug = "backup-site-west"
}

# Output comprehensive information showing relationships
output "site_and_group_integration" {
  description = "Demonstration of site and site group relationships"
  value = {
    created_site_group = {
      id          = data.netbox_site_group.production_group_data.id
      name        = data.netbox_site_group.production_group_data.name
      slug        = data.netbox_site_group.production_group_data.slug
      description = data.netbox_site_group.production_group_data.description
      tags        = data.netbox_site_group.production_group_data.tags
    }

    created_site = {
      id          = data.netbox_site.primary_dc_data.id
      name        = data.netbox_site.primary_dc_data.name
      slug        = data.netbox_site.primary_dc_data.slug
      status      = data.netbox_site.primary_dc_data.status
      group       = data.netbox_site.primary_dc_data.group
      description = data.netbox_site.primary_dc_data.description
      facility    = data.netbox_site.primary_dc_data.facility
    }

    existing_group = {
      id          = data.netbox_site_group.existing_dr_group.id
      name        = data.netbox_site_group.existing_dr_group.name
      slug        = data.netbox_site_group.existing_dr_group.slug
      description = data.netbox_site_group.existing_dr_group.description
    }

    existing_site = {
      id          = data.netbox_site.existing_site.id
      name        = data.netbox_site.existing_site.name
      slug        = data.netbox_site.existing_site.slug
      status      = data.netbox_site.existing_site.status
      group       = data.netbox_site.existing_site.group
      description = data.netbox_site.existing_site.description
    }
  }
}

# Example of validation and relationship checks
locals {
  site_belongs_to_group = data.netbox_site.primary_dc_data.group == data.netbox_site_group.production_group_data.name
  both_have_production_tag = (
    contains([for tag in data.netbox_site.primary_dc_data.tags : tag.name], "production") &&
    contains([for tag in data.netbox_site_group.production_group_data.tags : tag.name], "production")
  )
  existing_site_has_group = data.netbox_site.existing_site.group != null
}

output "relationship_validation" {
  value = {
    site_belongs_to_created_group = local.site_belongs_to_group
    consistent_tagging            = local.both_have_production_tag
    existing_site_has_group       = local.existing_site_has_group
    validation_status             = (local.site_belongs_to_group && local.both_have_production_tag) ? "All relationships valid" : "Check relationships"
  }
}

# Example of using data sources for resource dependencies
# This could be used to create additional sites in the same group
locals {
  production_group_id = data.netbox_site_group.production_group_data.id
}

# You could use this for creating additional resources
output "group_reference_for_new_sites" {
  description = "Group ID that can be used for creating additional sites"
  value = {
    group_id      = local.production_group_id
    group_name    = data.netbox_site_group.production_group_data.name
    usage_example = "Use group_id '${local.production_group_id}' in new site resources"
  }
}
