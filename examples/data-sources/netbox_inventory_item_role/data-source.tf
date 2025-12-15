data "netbox_inventory_item_role" "test" {
  name = "test-inventory-item-role"
}

output "example" {
  value = data.netbox_inventory_item_role.test.id
}
