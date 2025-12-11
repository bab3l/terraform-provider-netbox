# Look up a wireless link by ID
data "netbox_wireless_link" "example" {
  id = "123"
}

# Output the wireless link details
output "wireless_link_ssid" {
  value = data.netbox_wireless_link.example.ssid
}

output "wireless_link_status" {
  value = data.netbox_wireless_link.example.status
}

output "wireless_link_interface_a" {
  value = data.netbox_wireless_link.example.interface_a
}

output "wireless_link_interface_b" {
  value = data.netbox_wireless_link.example.interface_b
}
