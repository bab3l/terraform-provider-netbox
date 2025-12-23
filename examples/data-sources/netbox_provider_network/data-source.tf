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
