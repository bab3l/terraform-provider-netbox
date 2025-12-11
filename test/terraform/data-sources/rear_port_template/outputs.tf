output "id_matches" {
  value = data.netbox_rear_port_template.by_id.id == netbox_rear_port_template.test.id
}

output "name_matches" {
  value = data.netbox_rear_port_template.by_id.name == netbox_rear_port_template.test.name
}

output "positions_match" {
  value = data.netbox_rear_port_template.by_id.positions == netbox_rear_port_template.test.positions
}

output "type_matches" {
  value = data.netbox_rear_port_template.by_id.type == netbox_rear_port_template.test.type
}

output "description_matches" {
  value = data.netbox_rear_port_template.by_id.description == netbox_rear_port_template.test.description
}
