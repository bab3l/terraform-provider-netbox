output "source_id" {
  value = netbox_manufacturer.source.id
}

output "by_id_name" {
  value = data.netbox_manufacturer.by_id.name
}

output "by_name_id" {
  value = data.netbox_manufacturer.by_name.id
}

output "by_slug_id" {
  value = data.netbox_manufacturer.by_slug.id
}

output "all_ids_match" {
  value = data.netbox_manufacturer.by_id.id == data.netbox_manufacturer.by_name.id && data.netbox_manufacturer.by_name.id == data.netbox_manufacturer.by_slug.id
}
