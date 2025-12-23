data "netbox_inventory_item" "by_id" {
  id = 1
}

data "netbox_inventory_item" "by_name_device" {
  name      = "Module1"
  device_id = 5
}

output "item_id" {
  value = data.netbox_inventory_item.by_id.id
}

output "item_name" {
  value = data.netbox_inventory_item.by_name_device.name
}

output "item_device" {
  value = data.netbox_inventory_item.by_name_device.device_id
}
