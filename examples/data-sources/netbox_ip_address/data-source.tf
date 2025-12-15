data "netbox_ip_address" "test" {
  ip_address = "10.0.0.1/24"
}

output "example" {
  value = data.netbox_ip_address.test.id
}
