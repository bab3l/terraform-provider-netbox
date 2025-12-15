data "netbox_module_type" "test" {
  model = "test-module-type"
}

output "example" {
  value = data.netbox_module_type.test.id
}
