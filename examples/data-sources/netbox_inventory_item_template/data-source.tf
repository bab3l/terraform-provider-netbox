data "netbox_inventory_item_template" "by_id" {
  id = "1"
}

output "template_id" {
  value = data.netbox_inventory_item_template.by_id.id
}

output "template_name" {
  value = data.netbox_inventory_item_template.by_id.name
}

output "template_device_type" {
  value = data.netbox_inventory_item_template.by_id.device_type_id
}
