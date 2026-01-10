data "netbox_ipsec_policy" "by_id" {
  id = "1"
}

data "netbox_ipsec_policy" "by_name" {
  name = "Corporate-VPN-Policy"
}

output "policy_id" {
  value = data.netbox_ipsec_policy.by_id.id
}

output "policy_name" {
  value = data.netbox_ipsec_policy.by_name.name
}

output "policy_pfs_group" {
  value = data.netbox_ipsec_policy.by_name.pfs_group
}

output "policy_proposals" {
  value = data.netbox_ipsec_policy.by_name.proposals
}

output "policy_description" {
  value = data.netbox_ipsec_policy.by_name.description
}

# Access all custom fields
output "policy_custom_fields" {
  value       = data.netbox_ipsec_policy.by_id.custom_fields
  description = "All custom fields defined in NetBox for this IPsec policy"
}

# Access specific custom field by name
output "policy_rekey_time" {
  value       = try([for cf in data.netbox_ipsec_policy.by_id.custom_fields : cf.value if cf.name == "rekey_time_seconds"][0], null)
  description = "Example: accessing a numeric custom field for rekey time"
}

output "policy_dpd_enabled" {
  value       = try([for cf in data.netbox_ipsec_policy.by_id.custom_fields : cf.value if cf.name == "dpd_enabled"][0], null)
  description = "Example: accessing a boolean custom field for DPD status"
}
