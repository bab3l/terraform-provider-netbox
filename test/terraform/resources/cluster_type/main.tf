# Cluster Type Resource Test

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

# Test 1: Basic cluster type creation
resource "netbox_cluster_type" "basic" {
  name = "Test Cluster Type Basic"
  slug = "test-cluster-type-basic"
}

# Test 2: Cluster type with all optional fields
resource "netbox_cluster_type" "complete" {
  name        = "Test Cluster Type Complete"
  slug        = "test-cluster-type-complete"
  description = "A VMware vSphere cluster type for testing"
}

# Test 3: Different cluster type for variety
resource "netbox_cluster_type" "kubernetes" {
  name        = "Kubernetes"
  slug        = "kubernetes"
  description = "Kubernetes container orchestration platform"
}

# Test 4: Another cluster type
resource "netbox_cluster_type" "proxmox" {
  name        = "Proxmox VE"
  slug        = "proxmox-ve"
  description = "Proxmox Virtual Environment hypervisor"
}

# Outputs for verification
output "basic_id" {
  value = netbox_cluster_type.basic.id
}

output "complete_name" {
  value = netbox_cluster_type.complete.name
}

output "kubernetes_slug" {
  value = netbox_cluster_type.kubernetes.slug
}
