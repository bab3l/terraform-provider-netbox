output "basic_id_valid" {
  value = netbox_front_port.basic.id != ""
}

output "basic_rear_port_match" {
  value = netbox_front_port.basic.rear_port == netbox_rear_port.test.id
}

output "complete_id_valid" {
  value = netbox_front_port.complete.id != ""
}

output "complete_label_match" {
  value = netbox_front_port.complete.label == "Front Port 1"
}

output "complete_color_match" {
  value = netbox_front_port.complete.color == "aa1409"
}

output "complete_rear_port_position_match" {
  value = netbox_front_port.complete.rear_port_position == 2
}

output "complete_mark_connected" {
  value = netbox_front_port.complete.mark_connected == true
}

output "complete_rear_port_match" {
  value = netbox_front_port.complete.rear_port == netbox_rear_port.test.id
}
