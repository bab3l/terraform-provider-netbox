# FHRP Group Integration Test
# Tests the netbox_fhrp_group resource with basic and complete configurations

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

# Basic FHRP Group with only required fields (VRRP2)
resource "netbox_fhrp_group" "basic_vrrp2" {
  protocol = "vrrp2"
  group_id = 1
}

# VRRP3 FHRP Group
resource "netbox_fhrp_group" "vrrp3" {
  protocol    = "vrrp3"
  group_id    = 2
  name        = "VRRP3 Group"
  description = "VRRP version 3 group for integration testing"
}

# HSRP FHRP Group with authentication
resource "netbox_fhrp_group" "hsrp" {
  protocol    = "hsrp"
  group_id    = 10
  name        = "HSRP Standby Group"
  description = "HSRP group with plaintext authentication"
  auth_type   = "plaintext"
  auth_key    = "hsrpkey123"
  comments    = "Created by terraform integration test"
}

# CARP FHRP Group
resource "netbox_fhrp_group" "carp" {
  protocol    = "carp"
  group_id    = 5
  name        = "CARP Failover Group"
  description = "CARP group for BSD-style failover"
}

# GLBP FHRP Group with MD5 authentication
resource "netbox_fhrp_group" "glbp" {
  protocol    = "glbp"
  group_id    = 100
  name        = "GLBP Load Balancing Group"
  description = "Gateway Load Balancing Protocol group"
  auth_type   = "md5"
  auth_key    = "md5secret456"
}

# ClusterXL FHRP Group (Check Point)
resource "netbox_fhrp_group" "clusterxl" {
  protocol    = "clusterxl"
  group_id    = 1
  name        = "Check Point ClusterXL"
  description = "Check Point firewall cluster"
}

# Other protocol type
resource "netbox_fhrp_group" "other" {
  protocol    = "other"
  group_id    = 50
  name        = "Custom HA Protocol"
  description = "Custom or vendor-specific HA protocol"
}

# Complete FHRP Group with all optional fields
resource "netbox_fhrp_group" "complete" {
  protocol    = "vrrp2"
  group_id    = 254
  name        = "Complete FHRP Group"
  description = "FHRP group with all available fields"
  auth_type   = "md5"
  auth_key    = "fullsecret789"
  comments    = "This is a complete FHRP group configuration for integration testing"
}

# Outputs for verification
output "basic_vrrp2_id" {
  value = netbox_fhrp_group.basic_vrrp2.id
}

output "hsrp_id" {
  value = netbox_fhrp_group.hsrp.id
}

output "complete_id" {
  value = netbox_fhrp_group.complete.id
}
