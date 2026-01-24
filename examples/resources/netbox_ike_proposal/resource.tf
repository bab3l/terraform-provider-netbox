resource "netbox_ike_proposal" "test" {
  name                     = "Test IKE Proposal"
  authentication_method    = "pre-shared-keys"
  encryption_algorithm     = "aes-256-cbc"
  authentication_algorithm = "sha-256"
  group                    = "group-14"
}

import {
  to = netbox_ike_proposal.test
  id = "123"
}
