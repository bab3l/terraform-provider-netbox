# IPSec Proposal Resource Test

terraform {
  required_version = ">= 1.0"
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  # Uses NETBOX_SERVER_URL and NETBOX_API_TOKEN environment variables
}

# Test 1: Basic IPSec proposal with encryption algorithm
resource "netbox_ipsec_proposal" "basic" {
  name                 = "test-ipsec-proposal-basic"
  encryption_algorithm = "aes-256-cbc"
}

# Test 2: IPSec proposal with all optional fields
resource "netbox_ipsec_proposal" "complete" {
  name                     = "test-ipsec-proposal-complete"
  description              = "Complete IPSec proposal for integration testing"
  encryption_algorithm     = "aes-256-gcm"
  authentication_algorithm = "hmac-sha256"
  sa_lifetime_seconds      = 3600
  sa_lifetime_data         = 102400
  comments                 = "This IPSec proposal was created for integration testing."
}

# Test 3: IPSec proposal with only encryption
resource "netbox_ipsec_proposal" "encryption_only" {
  name                 = "test-ipsec-proposal-enc"
  encryption_algorithm = "aes-128-cbc"
}

# Test 4: Output values for verification
output "basic_proposal_id" {
  value = netbox_ipsec_proposal.basic.id
}

output "complete_proposal_name" {
  value = netbox_ipsec_proposal.complete.name
}

output "complete_proposal_encryption" {
  value = netbox_ipsec_proposal.complete.encryption_algorithm
}
