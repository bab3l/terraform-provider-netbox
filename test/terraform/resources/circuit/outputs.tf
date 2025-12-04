output "circuit_basic_id" {
  description = "ID of the basic circuit"
  value       = netbox_circuit.basic.id
}

output "circuit_basic_cid" {
  description = "CID of the basic circuit"
  value       = netbox_circuit.basic.cid
}

output "circuit_complete_id" {
  description = "ID of the complete circuit"
  value       = netbox_circuit.complete.id
}

output "circuit_complete_commit_rate" {
  description = "Commit rate of the complete circuit"
  value       = netbox_circuit.complete.commit_rate
}

output "circuit_active_status" {
  description = "Status of the active circuit"
  value       = netbox_circuit.active.status
}

output "basic_circuit_valid" {
  description = "Validates basic circuit was created correctly"
  value       = netbox_circuit.basic.id != "" && netbox_circuit.basic.cid == "CKT-BASIC-001"
}

output "complete_circuit_valid" {
  description = "Validates complete circuit was created correctly"
  value       = netbox_circuit.complete.id != "" && netbox_circuit.complete.commit_rate == 100000
}

output "provider_reference_valid" {
  description = "Validates circuit provider reference"
  value       = netbox_circuit.basic.circuit_provider == netbox_provider.test.slug
}

output "type_reference_valid" {
  description = "Validates circuit type reference"
  value       = netbox_circuit.basic.type == netbox_circuit_type.test.slug
}

output "tenant_association_valid" {
  description = "Validates circuit tenant association"
  value       = netbox_circuit.complete.tenant == netbox_tenant.test.id
}
