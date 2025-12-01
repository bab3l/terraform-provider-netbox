output "platform_id" {
  description = "ID of the created platform"
  value       = netbox_platform.basic.id
}

output "platform_manufacturer" {
  description = "Manufacturer ID referenced by the platform"
  value       = netbox_platform.basic.manufacturer
}
