data "netbox_wireless_lan" "test" {
  ssid = "test-ssid"
}

output "example" {
  value = data.netbox_wireless_lan.test.id
}
