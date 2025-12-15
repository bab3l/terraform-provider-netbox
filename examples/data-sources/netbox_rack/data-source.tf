data "netbox_rack" "test" {
  name = "test-rack"
}

output "example" {
  value = data.netbox_rack.test.id
}
