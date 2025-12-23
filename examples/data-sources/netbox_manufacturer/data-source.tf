data "netbox_manufacturer" "by_id" {
  id = "1"
}

data "netbox_manufacturer" "by_name" {
  name = "Cisco"
}

data "netbox_manufacturer" "by_slug" {
  slug = "cisco"
}

output "manufacturer_id" {
  value = data.netbox_manufacturer.by_id.id
}

output "manufacturer_name" {
  value = data.netbox_manufacturer.by_name.name
}

output "manufacturer_display" {
  value = data.netbox_manufacturer.by_slug.display_name
}
