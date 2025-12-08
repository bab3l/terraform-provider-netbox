# IPSec Policy Data Source Test

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
resource "netbox_ipsec_proposal" "test" {
  name                 = "test-ipsec-proposal-for-policy-ds"
  encryption_algorithm = "aes-256-cbc"
}

resource "netbox_ipsec_policy" "test" {
  name        = "test-ipsec-policy-ds"
  description = "Test IPSec policy for data source"
  proposals   = [netbox_ipsec_proposal.test.id]
  pfs_group   = 14
}

# Test: Lookup IPSec policy by ID
data "netbox_ipsec_policy" "by_id" {
  id = netbox_ipsec_policy.test.id
}

# Test: Lookup IPSec policy by name
data "netbox_ipsec_policy" "by_name" {
  name = netbox_ipsec_policy.test.name

  depends_on = [netbox_ipsec_policy.test]
}

# Output values for verification
output "by_id_name" {
  value = data.netbox_ipsec_policy.by_id.name
}

output "by_name_id" {
  value = data.netbox_ipsec_policy.by_name.id
}

output "by_id_pfs_group" {
  value = data.netbox_ipsec_policy.by_id.pfs_group
}
