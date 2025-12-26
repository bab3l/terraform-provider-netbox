# Look up a wireless link by ID
data "netbox_wireless_link" "by_id" {
  id = "1"
}

# Use wireless link data in outputs
output "link_info" {
  value = {
    id            = data.netbox_wireless_link.by_id.id
    interface_a   = data.netbox_wireless_link.by_id.interface_a
    interface_b   = data.netbox_wireless_link.by_id.interface_b
    ssid          = data.netbox_wireless_link.by_id.ssid
    status        = data.netbox_wireless_link.by_id.status
    auth_type     = data.netbox_wireless_link.by_id.auth_type
    auth_cipher   = data.netbox_wireless_link.by_id.auth_cipher
    distance      = data.netbox_wireless_link.by_id.distance
    distance_unit = data.netbox_wireless_link.by_id.distance_unit
  }
}

output "wireless_link_interfaces" {
  value = {
    interface_a = data.netbox_wireless_link.by_id.interface_a
    interface_b = data.netbox_wireless_link.by_id.interface_b
  }
}
