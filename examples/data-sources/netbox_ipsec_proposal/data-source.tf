data "netbox_ipsec_proposal" "test" {
  name = "test-ipsec-proposal"
}

output "example" {
  value = data.netbox_ipsec_proposal.test.id
}
