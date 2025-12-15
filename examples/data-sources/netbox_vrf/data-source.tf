data "netbox_vrf" "test" {
  name = "test-vrf"
}

output "example" {
  value = data.netbox_vrf.test.id
}
