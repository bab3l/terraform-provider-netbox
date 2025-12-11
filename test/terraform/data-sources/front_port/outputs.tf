output "id_matches" {
  value = data.netbox_front_port.by_id.id == netbox_front_port.test.id
}

output "name_matches" {
  value = data.netbox_front_port.by_id.name == netbox_front_port.test.name
}

output "rear_port_matches" {
  value = data.netbox_front_port.by_id.rear_port == netbox_rear_port.test.id
}

output "type_matches" {
  value = data.netbox_front_port.by_id.type == netbox_front_port.test.type
}

output "description_matches" {
  value = data.netbox_front_port.by_id.description == netbox_front_port.test.description
}
