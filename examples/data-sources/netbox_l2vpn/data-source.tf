data "netbox_l2vpn" "test" {
  name = "test-l2vpn"
}

output "example" {
  value = data.netbox_l2vpn.test.id
}
