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

resource "netbox_ipsec_proposal" "test" {
  name                     = "Test IPSec Proposal"
  encryption_algorithm     = "aes-256-cbc"
  authentication_algorithm = "sha-256"
}

resource "netbox_ipsec_policy" "test" {
  name      = "Test IPSec Policy"
  proposals = [netbox_ipsec_proposal.test.id]
}

resource "netbox_ipsec_profile" "test" {
  name         = "Test IPSec Profile"
  mode         = "esp"
  ike_policy   = netbox_ike_policy.test.id
  ipsec_policy = netbox_ipsec_policy.test.id
}
