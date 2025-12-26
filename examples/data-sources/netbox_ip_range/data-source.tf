data "netbox_ip_range" "by_id" {
  id = "1"
}

data "netbox_ip_range" "by_addresses" {
  start_address = "10.0.0.0"
  end_address   = "10.0.0.255"
}

output "range_id" {
  value = data.netbox_ip_range.by_id.id
}

output "range_start" {
  value = data.netbox_ip_range.by_addresses.start_address
}

output "range_end" {
  value = data.netbox_ip_range.by_addresses.end_address
}

output "range_size" {
  value = data.netbox_ip_range.by_addresses.size
}
