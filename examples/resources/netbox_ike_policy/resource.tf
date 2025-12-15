resource "netbox_ike_proposal" "test" {
  name                     = "Test IKE Proposal"
  authentication_method    = "pre-shared-keys"
  encryption_algorithm     = "aes-256-cbc"
  authentication_algorithm = "sha-256"
  group                    = "group-14"
}

resource "netbox_ike_policy" "test" {
  name      = "Test IKE Policy"
  version   = "v2"
  mode      = "main"
  proposals = [netbox_ike_proposal.test.id]
}
