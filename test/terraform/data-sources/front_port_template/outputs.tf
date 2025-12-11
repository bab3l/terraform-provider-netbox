output "id_matches" {
  value = data.netbox_front_port_template.by_id.id == netbox_front_port_template.test.id
}

output "name_matches" {
  value = data.netbox_front_port_template.by_id.name == netbox_front_port_template.test.name
}

output "rear_port_matches" {
  value = data.netbox_front_port_template.by_id.rear_port == netbox_rear_port_template.test.name
}

output "type_matches" {
  value = data.netbox_front_port_template.by_id.type == netbox_front_port_template.test.type
}

output "description_matches" {
  value = data.netbox_front_port_template.by_id.description == netbox_front_port_template.test.description
}
