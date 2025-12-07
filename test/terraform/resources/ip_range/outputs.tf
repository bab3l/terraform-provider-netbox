# Outputs to verify resource creation

output "basic_id" {
  description = "ID of the basic IP range"
  value       = netbox_ip_range.basic.id
}

output "basic_start_address" {
  description = "Start address of the basic IP range"
  value       = netbox_ip_range.basic.start_address
}

output "basic_end_address" {
  description = "End address of the basic IP range"
  value       = netbox_ip_range.basic.end_address
}

output "basic_id_valid" {
  description = "Basic IP range has valid ID"
  value       = netbox_ip_range.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete IP range"
  value       = netbox_ip_range.complete.id
}

output "complete_status" {
  description = "Status of the complete IP range"
  value       = netbox_ip_range.complete.status
}

output "complete_description" {
  description = "Description of the complete IP range"
  value       = netbox_ip_range.complete.description
}
