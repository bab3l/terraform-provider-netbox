output "source_id" {
  value = netbox_platform.source.id
}

output "by_id_name" {
  value = data.netbox_platform.by_id.name
}

output "by_name_id" {
  value = data.netbox_platform.by_name.id
}

output "by_slug_id" {
  value = data.netbox_platform.by_slug.id
}

output "manufacturer_id_match" {
  value = data.netbox_manufacturer.lookup.id == netbox_manufacturer.for_platform.id
}

output "all_ids_match" {
  value = data.netbox_platform.by_id.id == data.netbox_platform.by_name.id && data.netbox_platform.by_name.id == data.netbox_platform.by_slug.id
}
