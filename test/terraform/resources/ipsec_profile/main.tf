# IPSec Profile Resource Test

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

# Dependencies - create IKE and IPSec policies first
resource "netbox_ike_policy" "for_profile" {
  name    = "test-ike-policy-for-profile"
  version = 2
}

resource "netbox_ipsec_policy" "for_profile" {
  name = "test-ipsec-policy-for-profile"
}

# Test 1: Basic IPSec profile with required fields only
resource "netbox_ipsec_profile" "basic" {
  name         = "test-ipsec-profile-basic"
  mode         = "esp"
  ike_policy   = netbox_ike_policy.for_profile.id
  ipsec_policy = netbox_ipsec_policy.for_profile.id
}

# Test 2: IPSec profile with all optional fields
resource "netbox_ipsec_profile" "complete" {
  name         = "test-ipsec-profile-complete"
  description  = "Complete IPSec profile for integration testing"
  mode         = "ah"
  ike_policy   = netbox_ike_policy.for_profile.id
  ipsec_policy = netbox_ipsec_policy.for_profile.id
  comments     = "This IPSec profile was created for integration testing."
}

# Test 3: Output values for verification
output "basic_profile_id" {
  value = netbox_ipsec_profile.basic.id
}

output "complete_profile_name" {
  value = netbox_ipsec_profile.complete.name
}

output "complete_profile_mode" {
  value = netbox_ipsec_profile.complete.mode
}
