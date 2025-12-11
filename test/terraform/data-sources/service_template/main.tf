# Service Template Data Source Integration Test
# Tests the netbox_service_template data source for looking up existing service templates

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

# Create service templates to look up
resource "netbox_service_template" "test_http" {
  name        = "Test HTTP Service Template DS"
  protocol    = "tcp"
  ports       = [80, 8080]
  description = "HTTP service template for data source testing"
}

resource "netbox_service_template" "test_dns" {
  name        = "Test DNS Service Template DS"
  protocol    = "udp"
  ports       = [53]
  description = "DNS service template for data source testing"
}

# Look up service template by ID
data "netbox_service_template" "by_id" {
  id = netbox_service_template.test_http.id
}

# Look up service template by name
data "netbox_service_template" "by_name" {
  name = netbox_service_template.test_dns.name
}

# Outputs for verification
output "by_id_name" {
  value = data.netbox_service_template.by_id.name
}

output "by_id_protocol" {
  value = data.netbox_service_template.by_id.protocol
}

output "by_id_ports" {
  value = data.netbox_service_template.by_id.ports
}

output "by_name_id" {
  value = data.netbox_service_template.by_name.id
}

output "by_name_protocol" {
  value = data.netbox_service_template.by_name.protocol
}

output "by_name_description" {
  value = data.netbox_service_template.by_name.description
}

# Validation outputs
output "id_lookup_matches" {
  value = data.netbox_service_template.by_id.id == netbox_service_template.test_http.id
}

output "name_lookup_matches" {
  value = data.netbox_service_template.by_name.id == netbox_service_template.test_dns.id
}

output "all_lookups_valid" {
  value = (data.netbox_service_template.by_id.id == netbox_service_template.test_http.id) && (data.netbox_service_template.by_name.id == netbox_service_template.test_dns.id)
}
