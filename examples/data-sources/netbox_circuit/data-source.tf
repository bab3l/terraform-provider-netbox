data "netbox_circuit" "test" {
  cid = "123456789"
}

output "example" {
  value = data.netbox_circuit.test.id
}
