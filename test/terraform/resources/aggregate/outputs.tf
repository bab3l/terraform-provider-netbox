# Outputs to verify resource creation

output "basic_id" {
  description = "ID of the basic aggregate"
  value       = netbox_aggregate.basic.id
}

output "basic_prefix" {
  description = "Prefix of the basic aggregate"
  value       = netbox_aggregate.basic.prefix
}

output "basic_id_valid" {
  description = "Basic aggregate has valid ID"
  value       = netbox_aggregate.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete aggregate"
  value       = netbox_aggregate.complete.id
}

output "complete_prefix" {
  description = "Prefix of the complete aggregate"
  value       = netbox_aggregate.complete.prefix
}

output "complete_description" {
  description = "Description of the complete aggregate"
  value       = netbox_aggregate.complete.description
}

output "rir_id" {
  description = "ID of the RIR"
  value       = netbox_rir.test.id
}
