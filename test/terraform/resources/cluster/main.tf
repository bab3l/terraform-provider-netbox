# Cluster Resource Test

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

# Dependency: Cluster Type
resource "netbox_cluster_type" "vmware" {
  name        = "VMware vSphere (Cluster Test)"
  slug        = "vmware-vsphere-cluster-test"
  description = "VMware cluster type for cluster testing"
}

resource "netbox_cluster_type" "kubernetes" {
  name        = "Kubernetes (Cluster Test)"
  slug        = "kubernetes-cluster-test"
  description = "Kubernetes cluster type for cluster testing"
}

# Test 1: Basic cluster creation
resource "netbox_cluster" "basic" {
  name = "Test Cluster Basic"
  type = netbox_cluster_type.vmware.slug

  depends_on = [netbox_cluster_type.vmware]
}

# Test 2: Cluster with all optional fields
resource "netbox_cluster" "complete" {
  name        = "Test Cluster Complete"
  type        = netbox_cluster_type.vmware.slug
  status      = "active"
  description = "A complete test cluster with all fields"
  comments    = "This cluster was created for integration testing purposes."

  depends_on = [netbox_cluster_type.vmware]
}

# Test 3: Staging cluster
resource "netbox_cluster" "staging" {
  name        = "Test Cluster Staging"
  type        = netbox_cluster_type.vmware.slug
  status      = "staging"
  description = "A cluster in staging status"

  depends_on = [netbox_cluster_type.vmware]
}

# Test 4: Kubernetes cluster
resource "netbox_cluster" "k8s_production" {
  name        = "Production Kubernetes"
  type        = netbox_cluster_type.kubernetes.slug
  status      = "active"
  description = "Production Kubernetes cluster"
  comments    = "Handles production workloads"

  depends_on = [netbox_cluster_type.kubernetes]
}

# Test 5: Decommissioning cluster
resource "netbox_cluster" "decommissioning" {
  name        = "Legacy Cluster"
  type        = netbox_cluster_type.vmware.slug
  status      = "decommissioning"
  description = "Cluster being decommissioned"

  depends_on = [netbox_cluster_type.vmware]
}

# Outputs for verification
output "basic_id" {
  value = netbox_cluster.basic.id
}

output "complete_name" {
  value = netbox_cluster.complete.name
}

output "complete_type" {
  value = netbox_cluster.complete.type
}

output "staging_status" {
  value = netbox_cluster.staging.status
}

output "k8s_cluster_id" {
  value = netbox_cluster.k8s_production.id
}
