data "netbox_rack_role" "test" {
  name = "test-role"
}

output "example" {
  value = data.netbox_rack_role.test.id
}
