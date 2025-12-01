// Manufacturer Data Source Test

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

resource "netbox_manufacturer" "source" {
  name = "DS Test Manufacturer"
  slug = "ds-test-manufacturer"
}

data "netbox_manufacturer" "by_id" {
  id = netbox_manufacturer.source.id
}

data "netbox_manufacturer" "by_name" {
  name = netbox_manufacturer.source.name
}

data "netbox_manufacturer" "by_slug" {
  slug = netbox_manufacturer.source.slug
}
