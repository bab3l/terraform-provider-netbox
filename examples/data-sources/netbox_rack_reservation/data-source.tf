# Lookup by ID
data "netbox_rack_reservation" "by_id" {
  id = "123"
}

output "by_id" {
  value = data.netbox_rack_reservation.by_id.rack
}

output "by_id_units" {
  value = data.netbox_rack_reservation.by_id.units
}
