data "netbox_custom_field" "test" {
  name = "test_field"
}

output "example" {
  value = data.netbox_custom_field.test.id
}
