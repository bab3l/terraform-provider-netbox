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
