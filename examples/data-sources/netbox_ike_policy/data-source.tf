# Example: Look up an IKE policy by ID
data "netbox_ike_policy" "by_id" {
  id = "1"
}

# Example: Look up an IKE policy by name
data "netbox_ike_policy" "by_name" {
  name = "IKEv2-Policy"
}

# Example: Use IKE policy data in other resources
output "ike_policy_id" {
  value = data.netbox_ike_policy.by_id.id
}

output "ike_policy_name" {
  value = data.netbox_ike_policy.by_name.name
}

output "ike_policy_version" {
  value = data.netbox_ike_policy.by_name.version
}

output "ike_policy_mode" {
  value = data.netbox_ike_policy.by_name.mode
}

output "ike_policy_description" {
  value = data.netbox_ike_policy.by_name.description
}

output "ike_policy_proposals" {
  value = data.netbox_ike_policy.by_name.proposals
}

# Access all custom fields
output "ike_policy_custom_fields" {
  value       = data.netbox_ike_policy.by_id.custom_fields
  description = "All custom fields defined in NetBox for this IKE policy"
}

# Access specific custom field by name
output "ike_policy_dpd_action" {
  value       = try([for cf in data.netbox_ike_policy.by_id.custom_fields : cf.value if cf.name == "dpd_action"][0], null)
  description = "Example: accessing a select custom field for DPD action"
}

output "ike_policy_lifetime_seconds" {
  value       = try([for cf in data.netbox_ike_policy.by_id.custom_fields : cf.value if cf.name == "lifetime_seconds"][0], null)
  description = "Example: accessing a numeric custom field for lifetime"
}
