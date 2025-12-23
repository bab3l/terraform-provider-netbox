# Lookup by ID
data "netbox_prefix" "by_id" {
  id = "123"
}

output "by_id" {
  value = data.netbox_prefix.by_id.prefix
}

# Lookup by prefix CIDR
data "netbox_prefix" "by_cidr" {
  prefix = "10.0.0.0/24"
}

output "by_cidr" {
  value = data.netbox_prefix.by_cidr.status
}
