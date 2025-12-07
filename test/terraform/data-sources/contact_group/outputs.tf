# Contact Group Data Source Outputs

output "by_id_name" {
  value = data.netbox_contact_group.by_id.name
}

output "by_id_slug" {
  value = data.netbox_contact_group.by_id.slug
}

output "by_name_id" {
  value = data.netbox_contact_group.by_name.id
}

output "by_name_description" {
  value = data.netbox_contact_group.by_name.description
}
