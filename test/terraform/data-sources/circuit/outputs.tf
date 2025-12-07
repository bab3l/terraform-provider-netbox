# Circuit Data Source Outputs

output "by_id_cid" {
  value = data.netbox_circuit.by_id.cid
}

output "by_id_status" {
  value = data.netbox_circuit.by_id.status
}

output "by_cid_id" {
  value = data.netbox_circuit.by_cid.id
}

output "by_cid_description" {
  value = data.netbox_circuit.by_cid.description
}
