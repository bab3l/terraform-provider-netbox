data "netbox_circuit" "by_cid" {
  cid = "123456789"
}

data "netbox_circuit" "by_id" {
  id = "456"
}

output "by_cid" {
  value = data.netbox_circuit.by_cid.id
}

output "by_id" {
  value = data.netbox_circuit.by_id.cid
}

output "circuit_provider" {
  value = data.netbox_circuit.by_cid.provider
}

output "circuit_type" {
  value = data.netbox_circuit.by_cid.type
}

output "circuit_status" {
  value = data.netbox_circuit.by_cid.status
}

output "circuit_tenant" {
  value = data.netbox_circuit.by_cid.tenant
}

output "circuit_description" {
  value = data.netbox_circuit.by_cid.description
}

# Access all custom fields
output "circuit_custom_fields" {
  value       = data.netbox_circuit.by_cid.custom_fields
  description = "All custom fields defined in NetBox for this circuit"
}

# Access specific custom field by name
output "circuit_bandwidth" {
  value       = try([for cf in data.netbox_circuit.by_cid.custom_fields : cf.value if cf.name == "bandwidth_mbps"][0], null)
  description = "Example: accessing a numeric custom field for circuit bandwidth"
}

output "circuit_contract_id" {
  value       = try([for cf in data.netbox_circuit.by_cid.custom_fields : cf.value if cf.name == "contract_id"][0], null)
  description = "Example: accessing a text custom field for contract ID"
}
