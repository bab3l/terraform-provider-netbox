# IPSec Proposal Data Source Test

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
  name                     = "test-ipsec-proposal-ds"
  description              = "Test IPSec proposal for data source"
  encryption_algorithm     = "aes-256-gcm"
  authentication_algorithm = "hmac-sha256"
  sa_lifetime_seconds      = 3600
  sa_lifetime_data         = 102400
}

# Test: Lookup IPSec proposal by ID
data "netbox_ipsec_proposal" "by_id" {
  id = netbox_ipsec_proposal.test.id
}

# Test: Lookup IPSec proposal by name
data "netbox_ipsec_proposal" "by_name" {
  name = netbox_ipsec_proposal.test.name

  depends_on = [netbox_ipsec_proposal.test]
}

# Output values for verification
output "by_id_name" {
  value = data.netbox_ipsec_proposal.by_id.name
}

output "by_name_id" {
  value = data.netbox_ipsec_proposal.by_name.id
}

output "by_id_encryption" {
  value = data.netbox_ipsec_proposal.by_id.encryption_algorithm
}

output "by_name_sa_lifetime" {
  value = data.netbox_ipsec_proposal.by_name.sa_lifetime_seconds
}
