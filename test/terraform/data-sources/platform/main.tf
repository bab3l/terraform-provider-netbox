// Platform Data Source Test

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

# Create a manufacturer and platform to be looked up
resource "netbox_manufacturer" "for_platform" {
  name = "DS Platform Manufacturer"
  slug = "ds-platform-manufacturer"
}

resource "netbox_platform" "source" {
  name         = "DS Test Platform"
  slug         = "ds-test-platform"
  manufacturer = netbox_manufacturer.for_platform.id
}

data "netbox_platform" "by_id" {
  id = netbox_platform.source.id
}

data "netbox_platform" "by_name" {
  name = netbox_platform.source.name
}

data "netbox_platform" "by_slug" {
  slug = netbox_platform.source.slug
}

data "netbox_manufacturer" "lookup" {
  id = netbox_manufacturer.for_platform.id
}
