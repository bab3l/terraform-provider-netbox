resource "netbox_site_group" "example" {
  name        = "Example Site Group"
  slug        = "example-site-group"
  description = "An example site group created with Terraform"
}
