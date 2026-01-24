resource "netbox_ipsec_proposal" "test" {
  name                     = "Test IPSec Proposal"
  encryption_algorithm     = "aes-256-cbc"
  authentication_algorithm = "sha-256"
}

resource "netbox_ipsec_policy" "test" {
  name      = "Test IPSec Policy"
  proposals = [netbox_ipsec_proposal.test.id]
}

import {
  to = netbox_ipsec_policy.test
  id = "123"
}
