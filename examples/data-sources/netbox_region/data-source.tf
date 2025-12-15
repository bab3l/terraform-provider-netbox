data "netbox_region" "test" {
  name = "test-region"
}

output "example" {
  value = data.netbox_region.test.id
}
