data "netbox_inventory_item_role" "by_id" {
  id = "1"
}

data "netbox_inventory_item_role" "by_name" {
  name = "Power Supply"
}

data "netbox_inventory_item_role" "by_slug" {
  slug = "power-supply"
}

output "role_id" {
  value = data.netbox_inventory_item_role.by_id.id
}

output "role_name" {
  value = data.netbox_inventory_item_role.by_name.name
}

output "role_color" {
  value = data.netbox_inventory_item_role.by_slug.color
}
