output "manufacturer_id" {
  description = "ID of the created manufacturer"
  value       = netbox_manufacturer.basic.id
}

output "manufacturer_name" {
  description = "Name of the created manufacturer"
  value       = netbox_manufacturer.basic.name
}
