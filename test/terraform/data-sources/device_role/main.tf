# Device Role Data Source Test

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

# First create a device role that we can query via data sources
resource "netbox_device_role" "test_source" {
  name        = "Data Source Test Device Role"
  slug        = "data-source-test-device-role"
  color       = "e74c3c"
  vm_role     = true
  description = "Device role used to test data sources"
}

# Look up the device role by ID
data "netbox_device_role" "by_id" {
  id = netbox_device_role.test_source.id
}

# Look up the device role by slug
data "netbox_device_role" "by_slug" {
  slug = netbox_device_role.test_source.slug
  depends_on = [netbox_device_role.test_source]
}

# Look up the device role by name
data "netbox_device_role" "by_name" {
  name = netbox_device_role.test_source.name
  depends_on = [netbox_device_role.test_source]
}

# Outputs to verify data source values match the resource
output "source_id" {
  value = netbox_device_role.test_source.id
}

output "source_name" {
  value = netbox_device_role.test_source.name
}

output "source_slug" {
  value = netbox_device_role.test_source.slug
}

output "source_color" {
  value = netbox_device_role.test_source.color
}

output "source_vm_role" {
  value = netbox_device_role.test_source.vm_role
}

output "source_description" {
  value = netbox_device_role.test_source.description
}

# Verify lookup by ID
output "by_id_name" {
  value = data.netbox_device_role.by_id.name
}

output "by_id_slug" {
  value = data.netbox_device_role.by_id.slug
}

output "by_id_color" {
  value = data.netbox_device_role.by_id.color
}

output "by_id_vm_role" {
  value = data.netbox_device_role.by_id.vm_role
}

output "by_id_description" {
  value = data.netbox_device_role.by_id.description
}

# Verify lookup by slug
output "by_slug_id" {
  value = data.netbox_device_role.by_slug.id
}

output "by_slug_name" {
  value = data.netbox_device_role.by_slug.name
}

output "by_slug_color" {
  value = data.netbox_device_role.by_slug.color
}

output "by_slug_vm_role" {
  value = data.netbox_device_role.by_slug.vm_role
}

output "by_slug_description" {
  value = data.netbox_device_role.by_slug.description
}

# Verify lookup by name
output "by_name_id" {
  value = data.netbox_device_role.by_name.id
}

output "by_name_slug" {
  value = data.netbox_device_role.by_name.slug
}

output "by_name_color" {
  value = data.netbox_device_role.by_name.color
}

output "by_name_vm_role" {
  value = data.netbox_device_role.by_name.vm_role
}

output "by_name_description" {
  value = data.netbox_device_role.by_name.description
}

# Validate that all three lookups found the same device role
output "id_match" {
  value = data.netbox_device_role.by_id.id == data.netbox_device_role.by_slug.id && data.netbox_device_role.by_id.id == data.netbox_device_role.by_name.id
}

output "name_match" {
  value = data.netbox_device_role.by_id.name == data.netbox_device_role.by_slug.name && data.netbox_device_role.by_id.name == data.netbox_device_role.by_name.name
}

output "slug_match" {
  value = data.netbox_device_role.by_id.slug == data.netbox_device_role.by_slug.slug && data.netbox_device_role.by_id.slug == data.netbox_device_role.by_name.slug
}

output "vm_role_match" {
  value = data.netbox_device_role.by_id.vm_role == data.netbox_device_role.by_slug.vm_role && data.netbox_device_role.by_id.vm_role == data.netbox_device_role.by_name.vm_role
}
