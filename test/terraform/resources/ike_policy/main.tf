# IKE Policy Resource Test

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

# Dependencies - create IKE proposals first
resource "netbox_ike_proposal" "for_policy" {
  name                  = "test-ike-proposal-for-policy"
  authentication_method = "preshared-keys"
  encryption_algorithm  = "aes-256-cbc"
  group                 = 14
}

resource "netbox_ike_proposal" "for_policy_2" {
  name                  = "test-ike-proposal-for-policy-2"
  authentication_method = "preshared-keys"
  encryption_algorithm  = "aes-128-cbc"
  group                 = 5
}

# Test 1: Basic IKE policy with required fields only (IKEv2, no mode)
resource "netbox_ike_policy" "basic" {
  name    = "test-ike-policy-basic"
  version = 2
}

# Test 2: IKEv1 policy with mode
resource "netbox_ike_policy" "ikev1" {
  name    = "test-ike-policy-ikev1"
  version = 1
  mode    = "aggressive"
}

# Test 3: IKE policy with all optional fields
resource "netbox_ike_policy" "complete" {
  name          = "test-ike-policy-complete"
  description   = "Complete IKE policy for integration testing"
  version       = 2
  proposals     = [netbox_ike_proposal.for_policy.id, netbox_ike_proposal.for_policy_2.id]
  preshared_key = "supersecretkey123"
  comments      = "This IKE policy was created for integration testing."
}

# Test 4: Output values for verification
output "basic_policy_id" {
  value = netbox_ike_policy.basic.id
}

output "complete_policy_name" {
  value = netbox_ike_policy.complete.name
}

output "complete_policy_version" {
  value = netbox_ike_policy.complete.version
}
