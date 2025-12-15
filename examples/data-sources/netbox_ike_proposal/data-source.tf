data "netbox_ike_proposal" "test" {
  name = "test-ike-proposal"
}

output "example" {
  value = data.netbox_ike_proposal.test.id
}
