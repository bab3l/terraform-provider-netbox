data "netbox_rack_reservation" "test" {
  id = 123
}

output "example" {
  value = data.netbox_rack_reservation.test.id
}
