# Site Resource Test

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

# Test 1: Basic site creation
resource "netbox_site" "basic" {
  name        = "Test Site Basic"
  slug        = "test-site-basic"
  status      = "active"
  description = "A basic test site created by Terraform integration tests"
}

# Test 2: Site with all optional fields
resource "netbox_site" "complete" {
  name        = "Test Site Complete"
  slug        = "test-site-complete"
  status      = "planned"
  facility    = "DC-01"
  description = "A complete test site with all fields"
  comments    = "This site was created for integration testing purposes."
}

# Test 3: Site with site group reference (depends on site_group test running first or use ID)
# resource "netbox_site" "with_group" {
#   name   = "Test Site With Group"
#   slug   = "test-site-with-group"
#   group  = netbox_site_group.parent.id
# }
