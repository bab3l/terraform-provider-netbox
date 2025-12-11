output "basic_id_valid" {
  value = netbox_front_port_template.basic.id != ""
}

output "basic_rear_port_match" {
  value = netbox_front_port_template.basic.rear_port == netbox_rear_port_template.test.name
}

output "complete_id_valid" {
  value = netbox_front_port_template.complete.id != ""
}

output "complete_label_match" {
  value = netbox_front_port_template.complete.label == "Front Port Template 1"
}

output "complete_color_match" {
  value = netbox_front_port_template.complete.color == "aa1409"
}

output "complete_rear_port_position_match" {
  value = netbox_front_port_template.complete.rear_port_position == 2
}

output "complete_rear_port_match" {
  value = netbox_front_port_template.complete.rear_port == netbox_rear_port_template.test.name
}

output "complete_description_match" {
  value = netbox_front_port_template.complete.description == "Front port template for testing"
}
