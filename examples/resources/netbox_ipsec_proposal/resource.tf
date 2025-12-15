resource "netbox_ipsec_proposal" "test" {
  name                     = "Test IPSec Proposal"
  encryption_algorithm     = "aes-256-cbc"
  authentication_algorithm = "sha-256"
}
