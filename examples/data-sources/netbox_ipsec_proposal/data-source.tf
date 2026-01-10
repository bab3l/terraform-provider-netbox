data "netbox_ipsec_proposal" "by_id" {
  id = "1"
}

data "netbox_ipsec_proposal" "by_name" {
  name = "AES256-SHA256"
}

output "proposal_id" {
  value = data.netbox_ipsec_proposal.by_id.id
}

output "proposal_name" {
  value = data.netbox_ipsec_proposal.by_name.name
}

output "proposal_encryption" {
  value = data.netbox_ipsec_proposal.by_name.encryption_algorithm
}

output "proposal_authentication" {
  value = data.netbox_ipsec_proposal.by_name.authentication_algorithm
}

output "proposal_sa_lifetime" {
  value = data.netbox_ipsec_proposal.by_name.sa_lifetime_seconds
}

# Access all custom fields
output "proposal_custom_fields" {
  value       = data.netbox_ipsec_proposal.by_id.custom_fields
  description = "All custom fields defined in NetBox for this IPsec proposal"
}

# Access specific custom field by name
output "proposal_vendor" {
  value       = try([for cf in data.netbox_ipsec_proposal.by_id.custom_fields : cf.value if cf.name == "vendor_name"][0], null)
  description = "Example: accessing a text custom field for vendor"
}

output "proposal_is_legacy" {
  value       = try([for cf in data.netbox_ipsec_proposal.by_id.custom_fields : cf.value if cf.name == "is_legacy"][0], null)
  description = "Example: accessing a boolean custom field for legacy status"
}
