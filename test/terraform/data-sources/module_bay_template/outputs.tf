output "id_matches" {
  value = data.netbox_module_bay_template.by_id.id == netbox_module_bay_template.test.id
}

output "name_matches" {
  value = data.netbox_module_bay_template.by_id.name == netbox_module_bay_template.test.name
}

output "label_matches" {
  value = data.netbox_module_bay_template.by_id.label == netbox_module_bay_template.test.label
}

output "position_matches" {
  value = data.netbox_module_bay_template.by_id.position == netbox_module_bay_template.test.position
}

output "description_matches" {
  value = data.netbox_module_bay_template.by_id.description == netbox_module_bay_template.test.description
}
