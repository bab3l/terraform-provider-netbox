output "basic_id" {
  description = "ID of the basic rack reservation"
  value       = netbox_rack_reservation.basic.id
}

output "basic_units" {
  description = "Units of the basic rack reservation"
  value       = netbox_rack_reservation.basic.units
}

output "basic_id_valid" {
  description = "Basic rack reservation has valid ID"
  value       = netbox_rack_reservation.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete rack reservation"
  value       = netbox_rack_reservation.complete.id
}

output "complete_units_count" {
  description = "Number of units in complete reservation"
  value       = length(netbox_rack_reservation.complete.units)
}

output "complete_description" {
  description = "Description of the complete rack reservation"
  value       = netbox_rack_reservation.complete.description
}

output "rack_id" {
  description = "ID of the parent rack"
  value       = netbox_rack.test.id
}

output "site_id" {
  description = "ID of the parent site"
  value       = netbox_site.test.id
}
