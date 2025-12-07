# IP Address Data Source Outputs

output "by_id_address" {
  value = data.netbox_ip_address.by_id.address
}

output "by_id_status" {
  value = data.netbox_ip_address.by_id.status
}

output "by_address_id" {
  value = data.netbox_ip_address.by_address.id
}

output "by_address_description" {
  value = data.netbox_ip_address.by_address.description
}
