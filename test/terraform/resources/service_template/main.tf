# Service Template Resource Integration Test
# Tests the netbox_service_template resource with basic and complete configurations

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

# Supporting resource: Tag for testing tag associations
resource "netbox_tag" "service_template_test" {
  name        = "service-template-test"
  slug        = "service-template-test"
  color       = "4caf50"
  description = "Tag for service template integration tests"
}

# Test 1: Basic service template with only required fields
resource "netbox_service_template" "basic" {
  name  = "Test Service Template Basic"
  ports = [80]
}

# Test 2: Service template with protocol
resource "netbox_service_template" "with_protocol" {
  name     = "Test Service Template TCP"
  protocol = "tcp"
  ports    = [443]
}

# Test 3: Service template with multiple ports
resource "netbox_service_template" "multi_port" {
  name     = "Test Service Template Multi Port"
  protocol = "tcp"
  ports    = [8080, 8081, 8082]
}

# Test 4: Complete service template with all optional fields
resource "netbox_service_template" "complete" {
  name        = "Test Service Template Complete"
  protocol    = "udp"
  ports       = [53, 5353]
  description = "Complete service template for integration testing"
  comments    = "Created by Terraform integration test"

  tags = [
    {
      name  = netbox_tag.service_template_test.name
      slug  = netbox_tag.service_template_test.slug
      color = netbox_tag.service_template_test.color
    }
  ]
}

# Test 5: SCTP protocol service template
resource "netbox_service_template" "sctp" {
  name        = "Test Service Template SCTP"
  protocol    = "sctp"
  ports       = [3868]
  description = "SCTP service template (e.g., Diameter)"
}
