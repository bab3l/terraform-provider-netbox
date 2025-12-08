# IPSec Policy Resource Test

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

# Dependencies - create IPSec proposals first
resource "netbox_ipsec_proposal" "for_policy" {
  name                 = "test-ipsec-proposal-for-policy"
  encryption_algorithm = "aes-256-cbc"
}

resource "netbox_ipsec_proposal" "for_policy_2" {
  name                 = "test-ipsec-proposal-for-policy-2"
  encryption_algorithm = "aes-128-cbc"
}

# Test 1: Basic IPSec policy with required fields only
resource "netbox_ipsec_policy" "basic" {
  name = "test-ipsec-policy-basic"
}

# Test 2: IPSec policy with all optional fields
resource "netbox_ipsec_policy" "complete" {
  name        = "test-ipsec-policy-complete"
  description = "Complete IPSec policy for integration testing"
  proposals   = [netbox_ipsec_proposal.for_policy.id, netbox_ipsec_proposal.for_policy_2.id]
  pfs_group   = 14
  comments    = "This IPSec policy was created for integration testing."
}

# Test 3: IPSec policy with PFS group only
resource "netbox_ipsec_policy" "with_pfs" {
  name      = "test-ipsec-policy-pfs"
  pfs_group = 19
}

# Test 4: Output values for verification
output "basic_policy_id" {
  value = netbox_ipsec_policy.basic.id
}

output "complete_policy_name" {
  value = netbox_ipsec_policy.complete.name
}

output "complete_policy_pfs_group" {
  value = netbox_ipsec_policy.complete.pfs_group
}
