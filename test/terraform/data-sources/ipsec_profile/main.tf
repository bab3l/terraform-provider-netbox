# IPSec Profile Data Source Test

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

# Dependencies - create resources to test data sources
resource "netbox_ike_policy" "test" {
  name    = "test-ike-policy-for-profile-ds"
  version = 2
}

resource "netbox_ipsec_policy" "test" {
  name = "test-ipsec-policy-for-profile-ds"
}

resource "netbox_ipsec_profile" "test" {
  name         = "test-ipsec-profile-ds"
  description  = "Test IPSec profile for data source"
  mode         = "esp"
  ike_policy   = netbox_ike_policy.test.id
  ipsec_policy = netbox_ipsec_policy.test.id
}

# Test: Lookup IPSec profile by ID
data "netbox_ipsec_profile" "by_id" {
  id = netbox_ipsec_profile.test.id
}

# Test: Lookup IPSec profile by name
data "netbox_ipsec_profile" "by_name" {
  name = netbox_ipsec_profile.test.name

  depends_on = [netbox_ipsec_profile.test]
}

# Output values for verification
output "by_id_name" {
  value = data.netbox_ipsec_profile.by_id.name
}

output "by_name_id" {
  value = data.netbox_ipsec_profile.by_name.id
}

output "by_id_mode" {
  value = data.netbox_ipsec_profile.by_id.mode
}

output "by_name_ike_policy" {
  value = data.netbox_ipsec_profile.by_name.ike_policy
}
