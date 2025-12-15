data "netbox_circuit_type" "test" {
  slug = "test-circuit-type"
}

output "example" {
  value = data.netbox_circuit_type.test.id
}
