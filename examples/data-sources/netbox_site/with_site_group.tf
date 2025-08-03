# Site Data Source with Site Group Integration

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

# Create a site group for organization
resource "netbox_site_group" "production_sites" {
  name        = "Production Sites"
  slug        = "production-sites"
  description = "All production data centers and facilities"
  
  tags = [
    {
      name = "production"
      slug = "production"
    }
  ]
}

# Create a site that belongs to the site group
resource "netbox_site" "primary_dc" {
  name        = "Primary Datacenter"
  slug        = "primary-dc"
  status      = "active"
  description = "Main production facility"
  group       = netbox_site_group.production_sites.id
  
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

# Use data source to read the site back
data "netbox_site" "primary_dc_info" {
  id = netbox_site.primary_dc.id
}

# Use data source to read an existing site by slug
data "netbox_site" "existing_site" {
  slug = "dr-site-west"
}

# Output information about both sites
output "primary_datacenter" {
  description = "Information about the primary datacenter"
  value = {
    id          = data.netbox_site.primary_dc_info.id
    name        = data.netbox_site.primary_dc_info.name
    slug        = data.netbox_site.primary_dc_info.slug
    status      = data.netbox_site.primary_dc_info.status
    group       = data.netbox_site.primary_dc_info.group
    description = data.netbox_site.primary_dc_info.description
    tags        = data.netbox_site.primary_dc_info.tags
  }
}

output "existing_site_info" {
  description = "Information about an existing DR site"
  value = {
    id          = data.netbox_site.existing_site.id
    name        = data.netbox_site.existing_site.name
    slug        = data.netbox_site.existing_site.slug
    status      = data.netbox_site.existing_site.status
    group       = data.netbox_site.existing_site.group
    region      = data.netbox_site.existing_site.region
    description = data.netbox_site.existing_site.description
    facility    = data.netbox_site.existing_site.facility
  }
}

# Example of using site data for conditional logic
locals {
  site_is_active = data.netbox_site.primary_dc_info.status == "active"
  site_has_group = data.netbox_site.primary_dc_info.group != null
}

output "site_status" {
  value = {
    active    = local.site_is_active
    has_group = local.site_has_group
    message   = local.site_is_active ? "Site is ready for production use" : "Site needs attention"
  }
}
