# Example: Look up an IKE proposal by ID
data "netbox_ike_proposal" "by_id" {
  id = "1"
}

# Example: Look up an IKE proposal by name
data "netbox_ike_proposal" "by_name" {
  name = "AES-256-SHA"
}

# Example: Use IKE proposal data in other resources
output "ike_proposal_id" {
  value = data.netbox_ike_proposal.by_id.id
}

output "ike_proposal_name" {
  value = data.netbox_ike_proposal.by_name.name
}

output "ike_proposal_encryption" {
  value = data.netbox_ike_proposal.by_name.encryption_algorithm
}

output "ike_proposal_authentication" {
  value = data.netbox_ike_proposal.by_name.authentication_algorithm
}

output "ike_proposal_group" {
  value = data.netbox_ike_proposal.by_name.group
}

output "ike_proposal_sa_lifetime" {
  value = data.netbox_ike_proposal.by_name.sa_lifetime
}

# Access all custom fields
output "ike_proposal_custom_fields" {
  value       = data.netbox_ike_proposal.by_id.custom_fields
  description = "All custom fields defined in NetBox for this IKE proposal"
}

# Access specific custom field by name
output "ike_proposal_vendor" {
  value       = try([for cf in data.netbox_ike_proposal.by_id.custom_fields : cf.value if cf.name == "vendor_name"][0], null)
  description = "Example: accessing a text custom field for vendor"
}

output "ike_proposal_is_standard" {
  value       = try([for cf in data.netbox_ike_proposal.by_id.custom_fields : cf.value if cf.name == "is_standard"][0], null)
  description = "Example: accessing a boolean custom field for standard compliance"
}
