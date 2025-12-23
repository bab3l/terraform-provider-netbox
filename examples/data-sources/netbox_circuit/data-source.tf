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
