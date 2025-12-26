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
