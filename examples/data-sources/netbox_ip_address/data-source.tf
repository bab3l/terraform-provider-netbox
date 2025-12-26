data "netbox_ip_address" "by_id" {
  id = "1"
}

data "netbox_ip_address" "by_address" {
  address = "192.168.1.1/32"
}

output "address_id" {
  value = data.netbox_ip_address.by_id.id
}

output "address_value" {
  value = data.netbox_ip_address.by_address.address
}

output "address_status" {
  value = data.netbox_ip_address.by_address.status
}
