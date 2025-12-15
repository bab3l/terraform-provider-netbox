data "netbox_device_type" "test" {
  model = "test-model"
}

output "example" {
  value = data.netbox_device_type.test.id
}
