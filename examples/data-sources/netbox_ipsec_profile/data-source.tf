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

output "profile_ipsec_policy" {
  value = data.netbox_ipsec_profile.by_name.ipsec_policy
}

output "profile_description" {
  value = data.netbox_ipsec_profile.by_name.description
}

# Access all custom fields
output "profile_custom_fields" {
  value       = data.netbox_ipsec_profile.by_id.custom_fields
  description = "All custom fields defined in NetBox for this IPsec profile"
}

# Access specific custom field by name
output "profile_site_name" {
  value       = try([for cf in data.netbox_ipsec_profile.by_id.custom_fields : cf.value if cf.name == "site_name"][0], null)
  description = "Example: accessing a text custom field for site name"
}

output "profile_priority" {
  value       = try([for cf in data.netbox_ipsec_profile.by_id.custom_fields : cf.value if cf.name == "priority"][0], null)
  description = "Example: accessing a numeric custom field for priority"
}
