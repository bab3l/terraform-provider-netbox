output "id_matches" {
  value = data.netbox_rack_reservation.by_id.id == netbox_rack_reservation.test.id
}

output "units_match" {
  value = data.netbox_rack_reservation.by_id.units == netbox_rack_reservation.test.units
}

output "description_matches" {
  value = data.netbox_rack_reservation.by_id.description == netbox_rack_reservation.test.description
}
