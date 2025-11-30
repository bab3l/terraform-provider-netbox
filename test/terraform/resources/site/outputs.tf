# Outputs to verify resource creation

output "basic_site_id" {
  description = "ID of the basic site"
  value       = netbox_site.basic.id
}

output "basic_site_name" {
  description = "Name of the basic site"
  value       = netbox_site.basic.name
}

output "basic_site_slug" {
  description = "Slug of the basic site"
  value       = netbox_site.basic.slug
}

output "basic_site_status" {
  description = "Status of the basic site"
  value       = netbox_site.basic.status
}

output "complete_site_id" {
  description = "ID of the complete site"
  value       = netbox_site.complete.id
}

output "complete_site_name" {
  description = "Name of the complete site"
  value       = netbox_site.complete.name
}

output "complete_site_facility" {
  description = "Facility of the complete site"
  value       = netbox_site.complete.facility
}

output "complete_site_description" {
  description = "Description of the complete site"
  value       = netbox_site.complete.description
}
