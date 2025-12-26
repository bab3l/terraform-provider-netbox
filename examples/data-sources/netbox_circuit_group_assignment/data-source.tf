data "netbox_circuit_group_assignment" "example" {
  id = "456"
}

output "example" {
  value = data.netbox_circuit_group_assignment.example.circuit_cid
}
