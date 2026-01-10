# Look up IP address by ID
data "netbox_ip_address" "by_id" {
  id = "1"
}

# Look up IP address by address
data "netbox_ip_address" "by_address" {
  address = "192.168.1.1/32"
}

# Use IP address data in other resources
output "address_value" {
  value = data.netbox_ip_address.by_id.address
}

output "address_status" {
  value = data.netbox_ip_address.by_address.status
}

output "address_vrf" {
  value = data.netbox_ip_address.by_id.vrf
}

output "address_tenant" {
  value = data.netbox_ip_address.by_id.tenant
}

output "address_dns_name" {
  value = data.netbox_ip_address.by_address.dns_name
}

output "address_assigned_object_type" {
  value = data.netbox_ip_address.by_id.assigned_object_type
}

# Access all custom fields
output "address_custom_fields" {
  value       = data.netbox_ip_address.by_id.custom_fields
  description = "All custom fields defined in NetBox for this IP address"
}

# Access specific custom fields by name
output "address_hostname" {
  value       = try([for cf in data.netbox_ip_address.by_id.custom_fields : cf.value if cf.name == "hostname"][0], null)
  description = "Example: accessing a text custom field"
}

output "address_monitored" {
  value       = try([for cf in data.netbox_ip_address.by_id.custom_fields : cf.value if cf.name == "monitored"][0], null)
  description = "Example: accessing a boolean custom field"
}
