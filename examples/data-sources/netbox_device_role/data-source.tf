data "netbox_device_role" "test" {
  name = "test-role"
}

output "example" {
  value = data.netbox_device_role.test.id
}
