resource "netbox_site" "example" {
  name        = "Example Site"
  slug        = "example-site"
  status      = "active"
  description = "An example site created with Terraform"
  facility    = "DC01"
  comments    = "This is a sample site configuration"
}
