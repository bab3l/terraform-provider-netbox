# Comprehensive Site and Site Group Data Sources Example

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

# Create a hierarchical site group structure
resource "netbox_site_group" "global" {
  name        = "Global Infrastructure"
  slug        = "global-infrastructure"
  description = "Top-level site group for all infrastructure"

  tags = [
    {
      name = "global"
      slug = "global"
    }
  ]
}

resource "netbox_site_group" "regional" {
  name        = "North America"
  slug        = "north-america"
  parent      = netbox_site_group.global.id
  description = "North American sites"

  tags = [
    {
      name = "region"
      slug = "region"
    }
  ]
}

resource "netbox_site_group" "country" {
  name        = "United States"
  slug        = "united-states"
  parent      = netbox_site_group.regional.id
  description = "United States sites"

  tags = [
    {
      name = "country"
      slug = "country"
    }
  ]
}

# Create sites in different groups
resource "netbox_site" "headquarters" {
  name        = "Corporate Headquarters"
  slug        = "headquarters"
  status      = "active"
  group       = netbox_site_group.country.id
  description = "Main corporate facility"
  facility    = "HQ Building"

  tags = [
    {
      name = "headquarters"
      slug = "headquarters"
    }
  ]
}

resource "netbox_site" "datacenter" {
  name        = "Primary Datacenter"
  slug        = "primary-dc"
  status      = "active"
  group       = netbox_site_group.country.id
  description = "Primary data center facility"
  facility    = "DC-01"

  tags = [
    {
      name = "datacenter"
      slug = "datacenter"
    },
    {
      name = "production"
      slug = "production"
    }
  ]
}

# Use data sources to read all created resources
data "netbox_site_group" "global_data" {
  id = netbox_site_group.global.id
}

data "netbox_site_group" "regional_data" {
  slug = netbox_site_group.regional.slug
}

data "netbox_site_group" "country_data" {
  id = netbox_site_group.country.id
}

data "netbox_site" "hq_data" {
  slug = netbox_site.headquarters.slug
}

data "netbox_site" "dc_data" {
  id = netbox_site.datacenter.id
}

# Read some existing resources for comparison
data "netbox_site_group" "existing_group" {
  slug = "legacy-infrastructure"
}

data "netbox_site" "existing_site" {
  slug = "backup-facility"
}

# Generate comprehensive output showing the complete infrastructure
output "infrastructure_overview" {
  description = "Complete overview of site groups and sites using data sources"
  value = {
    hierarchy = {
      global_group = {
        id          = data.netbox_site_group.global_data.id
        name        = data.netbox_site_group.global_data.name
        slug        = data.netbox_site_group.global_data.slug
        parent      = data.netbox_site_group.global_data.parent
        description = data.netbox_site_group.global_data.description
        tags        = data.netbox_site_group.global_data.tags
      }

      regional_group = {
        id          = data.netbox_site_group.regional_data.id
        name        = data.netbox_site_group.regional_data.name
        slug        = data.netbox_site_group.regional_data.slug
        parent      = data.netbox_site_group.regional_data.parent
        description = data.netbox_site_group.regional_data.description
        tags        = data.netbox_site_group.regional_data.tags
      }

      country_group = {
        id          = data.netbox_site_group.country_data.id
        name        = data.netbox_site_group.country_data.name
        slug        = data.netbox_site_group.country_data.slug
        parent      = data.netbox_site_group.country_data.parent
        description = data.netbox_site_group.country_data.description
        tags        = data.netbox_site_group.country_data.tags
      }
    }

    sites = {
      headquarters = {
        id          = data.netbox_site.hq_data.id
        name        = data.netbox_site.hq_data.name
        slug        = data.netbox_site.hq_data.slug
        status      = data.netbox_site.hq_data.status
        group       = data.netbox_site.hq_data.group
        description = data.netbox_site.hq_data.description
        facility    = data.netbox_site.hq_data.facility
        tags        = data.netbox_site.hq_data.tags
      }

      datacenter = {
        id          = data.netbox_site.dc_data.id
        name        = data.netbox_site.dc_data.name
        slug        = data.netbox_site.dc_data.slug
        status      = data.netbox_site.dc_data.status
        group       = data.netbox_site.dc_data.group
        description = data.netbox_site.dc_data.description
        facility    = data.netbox_site.dc_data.facility
        tags        = data.netbox_site.dc_data.tags
      }
    }

    existing_resources = {
      existing_group = {
        id          = data.netbox_site_group.existing_group.id
        name        = data.netbox_site_group.existing_group.name
        slug        = data.netbox_site_group.existing_group.slug
        parent      = data.netbox_site_group.existing_group.parent
        description = data.netbox_site_group.existing_group.description
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
}

# Validate the hierarchical structure
locals {
  global_is_top_level     = data.netbox_site_group.global_data.parent == null
  regional_parent_correct = data.netbox_site_group.regional_data.parent == data.netbox_site_group.global_data.name
  country_parent_correct  = data.netbox_site_group.country_data.parent == data.netbox_site_group.regional_data.name

  hq_group_correct = data.netbox_site.hq_data.group == data.netbox_site_group.country_data.name
  dc_group_correct = data.netbox_site.dc_data.group == data.netbox_site_group.country_data.name

  hierarchy_valid = (
    local.global_is_top_level &&
    local.regional_parent_correct &&
    local.country_parent_correct &&
    local.hq_group_correct &&
    local.dc_group_correct
  )
}

output "hierarchy_validation" {
  description = "Validation of the site group and site hierarchy"
  value = {
    global_is_top_level        = local.global_is_top_level
    regional_hierarchy_correct = local.regional_parent_correct
    country_hierarchy_correct  = local.country_parent_correct
    hq_assignment_correct      = local.hq_group_correct
    dc_assignment_correct      = local.dc_group_correct
    overall_hierarchy_valid    = local.hierarchy_valid
    message                    = local.hierarchy_valid ? "All hierarchical relationships are correct" : "Check hierarchy configuration"
  }
}

# Example of using data sources for conditional resource creation
locals {
  country_group_id = data.netbox_site_group.country_data.id
  production_sites = [
    for site_key, site_data in {
      hq = data.netbox_site.hq_data
      dc = data.netbox_site.dc_data
      } : {
      name = site_data.name
      id   = site_data.id
      slug = site_data.slug
    }
    if site_data.status == "active"
  ]
}

output "operational_summary" {
  description = "Summary for operational use"
  value = {
    country_group_for_new_sites = local.country_group_id
    active_production_sites     = local.production_sites
    site_count_by_status = {
      active = length([
        for site in [data.netbox_site.hq_data, data.netbox_site.dc_data]
        : site if site.status == "active"
      ])
    }
    usage_note = "Use country_group_id '${local.country_group_id}' for new sites in the US"
  }
}
