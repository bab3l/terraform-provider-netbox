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
