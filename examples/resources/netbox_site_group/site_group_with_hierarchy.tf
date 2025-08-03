# Example Terraform configuration demonstrating hierarchical site groups with tags and custom fields
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

# Top-level site group for a geographic region
resource "netbox_site_group" "north_america" {
  name        = "North America"
  slug        = "north-america"
  description = "Site group for North American locations"

  # Tags for categorization
  tags = [
    {
      name = "geographic"
      slug = "geographic"
    },
    {
      name = "region"
      slug = "region"
    }
  ]

  # Custom fields for additional metadata
  custom_fields = [
    {
      name  = "region_code"
      type  = "text"
      value = "NA"
    },
    {
      name  = "time_zones"
      type  = "text"
      value = "EST,CST,MST,PST"
    }
  ]
}

# Child site group for a specific country
resource "netbox_site_group" "united_states" {
  name        = "United States"
  slug        = "united-states"
  parent      = netbox_site_group.north_america.id
  description = "Site group for United States locations"

  tags = [
    {
      name = "country"
      slug = "country"
    },
    {
      name = "usa"
      slug = "usa"
    }
  ]

  custom_fields = [
    {
      name  = "country_code"
      type  = "text"
      value = "US"
    },
    {
      name  = "currency"
      type  = "text"
      value = "USD"
    }
  ]
}

# Another child site group for states/provinces
resource "netbox_site_group" "california" {
  name        = "California"
  slug        = "california"
  parent      = netbox_site_group.united_states.id
  description = "Site group for California locations"

  tags = [
    {
      name = "state"
      slug = "state"
    },
    {
      name = "west-coast"
      slug = "west-coast"
    }
  ]

  custom_fields = [
    {
      name  = "state_code"
      type  = "text"
      value = "CA"
    },
    {
      name  = "tax_rate"
      type  = "text"
      value = "0.0925"
    }
  ]
}

# Example of a site that would use these site groups
resource "netbox_site" "san_francisco_dc" {
  name   = "San Francisco Data Center"
  slug   = "san-francisco-dc"
  group  = netbox_site_group.california.id
  status = "active"

  description = "Primary data center in San Francisco"
  facility    = "Building 1 - Floor 3"

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
