# Lookup by ID
data "netbox_provider_network" "by_id" {
  id = "123"
}

output "by_id" {
  value = data.netbox_provider_network.by_id.name
}

# Lookup by name
data "netbox_provider_network" "by_name" {
  name = "Primary Network"
}

output "by_name" {
  value = data.netbox_provider_network.by_name.service_id
}

# Lookup by name and circuit_provider filter
data "netbox_provider_network" "by_name_and_provider" {
  name             = "Primary Network"
  circuit_provider = "456"
}

output "by_name_and_provider" {
  value = data.netbox_provider_network.by_name_and_provider.description
}

output "network_provider" {
  value = data.netbox_provider_network.by_id.circuit_provider
}

output "network_service_id" {
  value = data.netbox_provider_network.by_id.service_id
}

output "network_description" {
  value = data.netbox_provider_network.by_id.description
}

# Access all custom fields
output "network_custom_fields" {
  value       = data.netbox_provider_network.by_id.custom_fields
  description = "All custom fields defined in NetBox for this provider network"
}

# Access specific custom field by name
output "network_vlan_range" {
  value       = try([for cf in data.netbox_provider_network.by_id.custom_fields : cf.value if cf.name == "vlan_range"][0], null)
  description = "Example: accessing a text custom field for VLAN range"
}

output "network_is_mpls" {
  value       = try([for cf in data.netbox_provider_network.by_id.custom_fields : cf.value if cf.name == "is_mpls"][0], null)
  description = "Example: accessing a boolean custom field for MPLS capability"
}
