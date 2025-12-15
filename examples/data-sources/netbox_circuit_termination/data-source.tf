data "netbox_circuit_termination" "test" {
  id = 123
}

output "example" {
  value = data.netbox_circuit_termination.test.id
}
