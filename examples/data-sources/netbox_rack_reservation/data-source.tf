# Lookup by ID
data "netbox_rack_reservation" "by_id" {
  id = "123"
}

# Use rack reservation data in other resources
output "reservation_rack" {
  value = data.netbox_rack_reservation.by_id.rack
}

output "reservation_units" {
  value = data.netbox_rack_reservation.by_id.units
}

output "reservation_description" {
  value = data.netbox_rack_reservation.by_id.description
}

output "reservation_tenant" {
  value = data.netbox_rack_reservation.by_id.tenant
}

output "reservation_user" {
  value = data.netbox_rack_reservation.by_id.user
}

# Access all custom fields
output "reservation_custom_fields" {
  value       = data.netbox_rack_reservation.by_id.custom_fields
  description = "All custom fields defined in NetBox for this rack reservation"
}

# Access specific custom fields by name
output "reservation_project_id" {
  value       = try([for cf in data.netbox_rack_reservation.by_id.custom_fields : cf.value if cf.name == "project_id"][0], null)
  description = "Example: accessing a text custom field"
}

output "reservation_end_date" {
  value       = try([for cf in data.netbox_rack_reservation.by_id.custom_fields : cf.value if cf.name == "end_date"][0], null)
  description = "Example: accessing a date custom field"
}
