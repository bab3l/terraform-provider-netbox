# IKE Policy Data Source Test

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

# Dependencies - create resources to test data sources
resource "netbox_ike_proposal" "test" {
  name                  = "test-ike-proposal-for-policy-ds"
  authentication_method = "preshared-keys"
  encryption_algorithm  = "aes-256-cbc"
  group                 = 14
}

resource "netbox_ike_policy" "test" {
  name        = "test-ike-policy-ds"
  description = "Test IKE policy for data source"
  version     = 2
  proposals   = [netbox_ike_proposal.test.id]
}

# Test: Lookup IKE policy by ID
data "netbox_ike_policy" "by_id" {
  id = netbox_ike_policy.test.id
}

# Test: Lookup IKE policy by name
data "netbox_ike_policy" "by_name" {
  name = netbox_ike_policy.test.name

  depends_on = [netbox_ike_policy.test]
}

# Output values for verification
output "by_id_name" {
  value = data.netbox_ike_policy.by_id.name
}

output "by_name_id" {
  value = data.netbox_ike_policy.by_name.id
}

output "by_id_version" {
  value = data.netbox_ike_policy.by_id.version
}
