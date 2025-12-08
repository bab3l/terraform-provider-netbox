# IKE Proposal Resource Test

terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  # Uses NETBOX_SERVER_URL and NETBOX_API_TOKEN environment variables
}

# Test 1: Basic IKE proposal with required fields only
resource "netbox_ike_proposal" "basic" {
  name                  = "test-ike-proposal-basic"
  authentication_method = "preshared-keys"
  encryption_algorithm  = "aes-256-cbc"
  group                 = 14
}

# Test 2: IKE proposal with all optional fields
resource "netbox_ike_proposal" "complete" {
  name                     = "test-ike-proposal-complete"
  description              = "Complete IKE proposal for integration testing"
  authentication_method    = "certificates"
  encryption_algorithm     = "aes-256-gcm"
  authentication_algorithm = "hmac-sha256"
  group                    = 19
  sa_lifetime              = 28800
  comments                 = "This IKE proposal was created for integration testing."
}

# Test 3: Output values for verification
output "basic_proposal_id" {
  value = netbox_ike_proposal.basic.id
}

output "complete_proposal_name" {
  value = netbox_ike_proposal.complete.name
}

output "complete_proposal_encryption" {
  value = netbox_ike_proposal.complete.encryption_algorithm
}
