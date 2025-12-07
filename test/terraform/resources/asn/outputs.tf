# Outputs to verify resource creation

output "basic_id" {
  description = "ID of the basic ASN"
  value       = netbox_asn.basic.id
}

output "basic_asn" {
  description = "ASN number of the basic ASN"
  value       = netbox_asn.basic.asn
}

output "basic_id_valid" {
  description = "Basic ASN has valid ID"
  value       = netbox_asn.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete ASN"
  value       = netbox_asn.complete.id
}

output "complete_asn" {
  description = "ASN number of the complete ASN"
  value       = netbox_asn.complete.asn
}

output "complete_description" {
  description = "Description of the complete ASN"
  value       = netbox_asn.complete.description
}

output "rir_id" {
  description = "ID of the RIR"
  value       = netbox_rir.test.id
}
