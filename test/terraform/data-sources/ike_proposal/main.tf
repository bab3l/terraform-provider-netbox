# IKE Proposal Data Source Test

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
resource "netbox_ike_proposal" "test" {
  name                     = "test-ike-proposal-ds"
  description              = "Test IKE proposal for data source"
  authentication_method    = "preshared-keys"
  encryption_algorithm     = "aes-256-cbc"
  authentication_algorithm = "hmac-sha256"
  group                    = 14
  sa_lifetime              = 28800
}

# Test: Lookup IKE proposal by ID
data "netbox_ike_proposal" "by_id" {
  id = netbox_ike_proposal.test.id
}

# Test: Lookup IKE proposal by name
data "netbox_ike_proposal" "by_name" {
  name = netbox_ike_proposal.test.name

  depends_on = [netbox_ike_proposal.test]
}

# Output values for verification
output "by_id_name" {
  value = data.netbox_ike_proposal.by_id.name
}

output "by_name_id" {
  value = data.netbox_ike_proposal.by_name.id
}

output "by_id_encryption" {
  value = data.netbox_ike_proposal.by_id.encryption_algorithm
}

output "by_name_group" {
  value = data.netbox_ike_proposal.by_name.group
}
