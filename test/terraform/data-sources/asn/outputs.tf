# ASN Data Source Outputs

output "by_id_asn" {
  value = data.netbox_asn.by_id.asn
}

output "by_id_rir" {
  value = data.netbox_asn.by_id.rir
}

output "by_asn_id" {
  value = data.netbox_asn.by_asn.id
}

output "by_asn_description" {
  value = data.netbox_asn.by_asn.description
}
