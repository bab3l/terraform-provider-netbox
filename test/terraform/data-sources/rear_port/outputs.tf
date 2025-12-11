output "id_matches" {
  value = tostring(data.netbox_rear_port.by_id.id) == netbox_rear_port.test.id
}

output "name_matches" {
  value = data.netbox_rear_port.by_id.name == netbox_rear_port.test.name
}

output "type_matches" {
  value = data.netbox_rear_port.by_id.type == netbox_rear_port.test.type
}

output "positions_match" {
  value = data.netbox_rear_port.by_id.positions == netbox_rear_port.test.positions
}

output "description_matches" {
  value = data.netbox_rear_port.by_id.description == netbox_rear_port.test.description
}
