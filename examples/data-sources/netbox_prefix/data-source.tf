data "netbox_prefix" "test" {
  prefix = "10.0.0.0/24"
}

output "example" {
  value = data.netbox_prefix.test.id
}
