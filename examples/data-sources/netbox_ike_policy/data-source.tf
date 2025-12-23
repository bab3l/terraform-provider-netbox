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
