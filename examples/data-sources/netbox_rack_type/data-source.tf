data "netbox_rack_type" "test" {
  slug = "test-rack-type"
}

output "example" {
  value = data.netbox_rack_type.test.id
}
