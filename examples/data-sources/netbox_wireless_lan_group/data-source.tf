data "netbox_wireless_lan_group" "test" {
  name = "test-wireless-lan-group"
}

output "example" {
  value = data.netbox_wireless_lan_group.test.id
}
