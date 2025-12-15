data "netbox_inventory_item" "test" {
  name      = "test-inventory-item"
  device_id = 123
}

output "example" {
  value = data.netbox_inventory_item.test.id
}
