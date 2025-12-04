# Look up a cable by its ID
# Cables don't have names or slugs, so ID is the only lookup method

data "netbox_cable" "example" {
  id = "123"
}

# Use the retrieved cable data
output "cable_type" {
  value = data.netbox_cable.example.type
}

output "cable_status" {
  value = data.netbox_cable.example.status
}

output "cable_label" {
  value = data.netbox_cable.example.label
}

output "cable_length" {
  value = "${data.netbox_cable.example.length} ${data.netbox_cable.example.length_unit}"
}
