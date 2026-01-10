# Example for the plural/query IP addresses data source.
#
# Notes:
# - Multiple `filter` blocks are ANDed together.
# - Multiple values inside one filter block are ORed together.
# - The datasource returns `ids`, `addresses`, and `ip_addresses` (list of `{id,address}` objects).

resource "netbox_ip_address" "example" {
  address = "192.0.2.10/32"
  status  = "active"
}

data "netbox_ip_addresses" "by_address" {
  filter {
    name   = "address"
    values = [netbox_ip_address.example.address]
  }
}

output "ip_address_ids" {
  value = data.netbox_ip_addresses.by_address.ids
}

output "ip_address_addresses" {
  value = data.netbox_ip_addresses.by_address.addresses
}

output "ip_address_objects" {
  value = data.netbox_ip_addresses.by_address.ip_addresses
}
