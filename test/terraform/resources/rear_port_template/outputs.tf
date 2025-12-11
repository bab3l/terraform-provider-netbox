output "basic_id_valid" {
  value = netbox_rear_port_template.basic.id != ""
}

output "complete_id_valid" {
  value = netbox_rear_port_template.complete.id != ""
}

output "complete_label_match" {
  value = netbox_rear_port_template.complete.label == "Rear Port Template 1"
}

output "complete_color_match" {
  value = netbox_rear_port_template.complete.color == "aa1409"
}

output "complete_positions_match" {
  value = netbox_rear_port_template.complete.positions == 4
}

output "complete_description_match" {
  value = netbox_rear_port_template.complete.description == "Rear port template for testing"
}
