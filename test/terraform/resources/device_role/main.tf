# Device Role Resource Test

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

# Basic device role with only required fields
resource "netbox_device_role" "basic" {
  name = "Basic Test Device Role"
  slug = "basic-test-device-role"
}

# Complete device role with all optional fields
resource "netbox_device_role" "complete" {
  name        = "Complete Test Device Role"
  slug        = "complete-test-device-role"
  color       = "aa1409"
  vm_role     = true
  description = "This is a complete test device role with all fields."
}

# Device role for routers (vm_role = false since routers aren't VMs)
resource "netbox_device_role" "router" {
  name        = "Router Device Role"
  slug        = "router-device-role"
  color       = "2ecc71"
  vm_role     = false
  description = "Devices functioning as network routers."
}

# Device role for switches
resource "netbox_device_role" "switch" {
  name        = "Switch Device Role"
  slug        = "switch-device-role"
  color       = "f39c12"
  vm_role     = false
  description = "Network switches in the infrastructure."
}

# Device role for servers (can be VMs)
resource "netbox_device_role" "server" {
  name        = "Server Device Role"
  slug        = "server-device-role"
  color       = "3498db"
  vm_role     = true
  description = "Physical or virtual servers."
}

# Device role for firewalls
resource "netbox_device_role" "firewall" {
  name        = "Firewall Device Role"
  slug        = "firewall-device-role"
  color       = "e74c3c"
  vm_role     = true
  description = "Physical or virtual firewalls."
}

# Output values for verification
output "basic_id" {
  value = netbox_device_role.basic.id
}

output "basic_name" {
  value = netbox_device_role.basic.name
}

output "basic_slug" {
  value = netbox_device_role.basic.slug
}

output "basic_vm_role" {
  value = netbox_device_role.basic.vm_role
}

output "complete_id" {
  value = netbox_device_role.complete.id
}

output "complete_name" {
  value = netbox_device_role.complete.name
}

output "complete_slug" {
  value = netbox_device_role.complete.slug
}

output "complete_color" {
  value = netbox_device_role.complete.color
}

output "complete_vm_role" {
  value = netbox_device_role.complete.vm_role
}

output "complete_description" {
  value = netbox_device_role.complete.description
}

output "router_id" {
  value = netbox_device_role.router.id
}

output "router_color" {
  value = netbox_device_role.router.color
}

output "router_vm_role" {
  value = netbox_device_role.router.vm_role
}

output "switch_id" {
  value = netbox_device_role.switch.id
}

output "switch_color" {
  value = netbox_device_role.switch.color
}

output "switch_vm_role" {
  value = netbox_device_role.switch.vm_role
}

output "server_id" {
  value = netbox_device_role.server.id
}

output "server_color" {
  value = netbox_device_role.server.color
}

output "server_vm_role" {
  value = netbox_device_role.server.vm_role
}

output "firewall_id" {
  value = netbox_device_role.firewall.id
}

output "firewall_color" {
  value = netbox_device_role.firewall.color
}

output "firewall_vm_role" {
  value = netbox_device_role.firewall.vm_role
}
