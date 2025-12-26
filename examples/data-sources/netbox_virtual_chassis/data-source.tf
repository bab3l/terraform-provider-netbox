# Look up a virtual chassis by ID
data "netbox_virtual_chassis" "by_id" {
  id = "1"
}

# Look up a virtual chassis by name
data "netbox_virtual_chassis" "by_name" {
  name = "test-virtual-chassis"
}

# Use virtual chassis data in outputs
output "by_id" {
  value = data.netbox_virtual_chassis.by_id.name
}

output "by_name" {
  value = data.netbox_virtual_chassis.by_name.id
}

output "master_device" {
  value = data.netbox_virtual_chassis.by_name.master
}
