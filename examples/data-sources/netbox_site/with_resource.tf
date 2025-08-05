# Site Data Source with Resource Integration

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

# Create a new site
resource "netbox_site" "primary_datacenter" {
  name        = "Primary Datacenter"
  slug        = "primary-dc"
  status      = "active"
  description = "Our main data center facility"
  facility    = "Building A"
  comments    = "Primary site for all production workloads"

  tags = [
    {
      name = "production"
      slug = "production"
    },
    {
      name = "datacenter"
      slug = "datacenter"
    }
  ]

  custom_fields = [
    {
      name  = "cost_center"
      type  = "text"
      value = "IT-OPS-001"
    },
    {
      name  = "contact_email"
      type  = "text"
      value = "datacenter-ops@example.com"
    }
  ]
}

# Use a data source to read the created site by ID
data "netbox_site" "primary_datacenter_data" {
  id = netbox_site.primary_datacenter.id
}

# Use a data source to read the created site by slug
data "netbox_site" "primary_datacenter_by_slug" {
  slug = netbox_site.primary_datacenter.slug
}

# Output comprehensive site information
output "primary_datacenter_info" {
  description = "Complete information about the primary datacenter"
  value = {
    resource_data = {
      id            = netbox_site.primary_datacenter.id
      name          = netbox_site.primary_datacenter.name
      slug          = netbox_site.primary_datacenter.slug
      status        = netbox_site.primary_datacenter.status
      description   = netbox_site.primary_datacenter.description
      facility      = netbox_site.primary_datacenter.facility
      comments      = netbox_site.primary_datacenter.comments
      tags          = netbox_site.primary_datacenter.tags
      custom_fields = netbox_site.primary_datacenter.custom_fields
    }

    data_source_by_id = {
      id            = data.netbox_site.primary_datacenter_data.id
      name          = data.netbox_site.primary_datacenter_data.name
      slug          = data.netbox_site.primary_datacenter_data.slug
      status        = data.netbox_site.primary_datacenter_data.status
      description   = data.netbox_site.primary_datacenter_data.description
      facility      = data.netbox_site.primary_datacenter_data.facility
      comments      = data.netbox_site.primary_datacenter_data.comments
      tags          = data.netbox_site.primary_datacenter_data.tags
      custom_fields = data.netbox_site.primary_datacenter_data.custom_fields
    }

    data_source_by_slug = {
      id            = data.netbox_site.primary_datacenter_by_slug.id
      name          = data.netbox_site.primary_datacenter_by_slug.name
      slug          = data.netbox_site.primary_datacenter_by_slug.slug
      status        = data.netbox_site.primary_datacenter_by_slug.status
      description   = data.netbox_site.primary_datacenter_by_slug.description
      facility      = data.netbox_site.primary_datacenter_by_slug.facility
      comments      = data.netbox_site.primary_datacenter_by_slug.comments
      tags          = data.netbox_site.primary_datacenter_by_slug.tags
      custom_fields = data.netbox_site.primary_datacenter_by_slug.custom_fields
    }
  }
}
