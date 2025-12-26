data "netbox_ipsec_profile" "by_id" {
  id = "1"
}

data "netbox_ipsec_profile" "by_name" {
  name = "Site-to-Site-VPN"
}

output "profile_id" {
  value = data.netbox_ipsec_profile.by_id.id
}

output "profile_name" {
  value = data.netbox_ipsec_profile.by_name.name
}

output "profile_mode" {
  value = data.netbox_ipsec_profile.by_name.mode
}

output "profile_ike_policy" {
  value = data.netbox_ipsec_profile.by_name.ike_policy
}
