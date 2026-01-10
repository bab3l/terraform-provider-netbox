data "netbox_l2vpn" "by_id" {
  id = "1"
}

data "netbox_l2vpn" "by_name" {
  name = "Corporate-L2VPN"
}

data "netbox_l2vpn" "by_slug" {
  slug = "corporate-l2vpn"
}

output "l2vpn_id" {
  value = data.netbox_l2vpn.by_id.id
}

output "l2vpn_name" {
  value = data.netbox_l2vpn.by_name.name
}

output "l2vpn_slug" {
  value = data.netbox_l2vpn.by_slug.slug
}

output "l2vpn_type" {
  value = data.netbox_l2vpn.by_slug.type
}

output "l2vpn_identifier" {
  value = data.netbox_l2vpn.by_id.identifier
}

output "l2vpn_tenant" {
  value = data.netbox_l2vpn.by_id.tenant
}

output "l2vpn_description" {
  value = data.netbox_l2vpn.by_id.description
}

# Access all custom fields
output "l2vpn_custom_fields" {
  value       = data.netbox_l2vpn.by_id.custom_fields
  description = "All custom fields defined in NetBox for this L2VPN"
}

# Access specific custom field by name
output "l2vpn_vlan_id" {
  value       = try([for cf in data.netbox_l2vpn.by_id.custom_fields : cf.value if cf.name == "vlan_id"][0], null)
  description = "Example: accessing a numeric custom field for VLAN ID"
}

output "l2vpn_encryption_type" {
  value       = try([for cf in data.netbox_l2vpn.by_id.custom_fields : cf.value if cf.name == "encryption_type"][0], null)
  description = "Example: accessing a select custom field for encryption type"
}
