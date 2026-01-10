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

# Access all custom fields
output "cable_custom_fields" {
  value       = data.netbox_cable.example.custom_fields
  description = "All custom fields defined in NetBox for this cable"
}

# Access a specific custom field by name
output "cable_vendor_id" {
  value       = try([for cf in data.netbox_cable.example.custom_fields : cf.value if cf.name == "vendor_id"][0], null)
  description = "Example: accessing a specific custom field value"
}
