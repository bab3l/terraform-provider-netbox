data "netbox_role" "test" {
  name = "test-role"
}

output "example" {
  value = data.netbox_role.test.id
}
