data "netbox_manufacturer" "test" {
  name = "test-manufacturer"
}

output "example" {
  value = data.netbox_manufacturer.test.id
}
