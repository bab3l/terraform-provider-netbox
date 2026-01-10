# Example for the plural/query interfaces data source.
#
# Notes:
# - Multiple `filter` blocks are ANDed together.
# - Multiple values inside one filter block are ORed together.
# - The datasource returns `ids`, `names`, and `interfaces` (list of `{id,name}` objects).

resource "netbox_site" "example" {
  name = "example"
  slug = "example"
}

resource "netbox_device_role" "example" {
  name = "example"
  slug = "example"
}

resource "netbox_manufacturer" "example" {
  name = "example"
  slug = "example"
}

resource "netbox_device_type" "example" {
  manufacturer = netbox_manufacturer.example.id
  model        = "example"
  slug         = "example"
}

resource "netbox_device" "example" {
  name        = "example"
  device_type = netbox_device_type.example.id
  role        = netbox_device_role.example.id
  site        = netbox_site.example.id
}

resource "netbox_interface" "example" {
  device = netbox_device.example.id
  name   = "eth0"
  type   = "1000base-t"
}

data "netbox_interfaces" "by_device_and_name" {
  filter {
    name   = "device_id"
    values = [netbox_device.example.id]
  }

  filter {
    name   = "name"
    values = [netbox_interface.example.name]
  }
}

output "interface_ids" {
  value = data.netbox_interfaces.by_device_and_name.ids
}

output "interface_names" {
  value = data.netbox_interfaces.by_device_and_name.names
}

output "interface_objects" {
  value = data.netbox_interfaces.by_device_and_name.interfaces
}
