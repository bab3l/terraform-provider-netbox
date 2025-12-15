data "netbox_inventory_item_template" "test" {
  name           = "test-inventory-item-template"
  device_type_id = 123
}

output "example" {
  value = data.netbox_inventory_item_template.test.id
}
