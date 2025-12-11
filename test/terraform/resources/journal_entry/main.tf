# Journal Entry Integration Test
# Tests the netbox_journal_entry resource with basic and complete configurations

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

# Create a site to attach journal entries to
resource "netbox_site" "test" {
  name        = "Journal Test Site"
  slug        = "journal-test-site"
  description = "Site for journal entry testing"
}

# Basic journal entry with only required fields
resource "netbox_journal_entry" "basic" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "Basic journal entry for testing"
}

# Journal entry with info kind (default)
resource "netbox_journal_entry" "info" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "This is an informational note about the site."
  kind                 = "info"
}

# Journal entry with success kind
resource "netbox_journal_entry" "success" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "Successfully completed site configuration!"
  kind                 = "success"
}

# Journal entry with warning kind
resource "netbox_journal_entry" "warning" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "Warning: This site needs attention - capacity at 80%"
  kind                 = "warning"
}

# Journal entry with danger kind
resource "netbox_journal_entry" "danger" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "CRITICAL: Site outage detected - immediate action required!"
  kind                 = "danger"
}

# Journal entry with markdown content
resource "netbox_journal_entry" "markdown" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  kind                 = "info"
  comments             = "# Site Documentation\n\n## Overview\nThis site has been configured for production use.\n\n## Key Points\n- Network configured\n- Power redundancy enabled\n- Monitoring active\n\n## Contacts\n- Primary: admin@example.com\n- Backup: support@example.com"
}

# Create a device for additional testing
resource "netbox_manufacturer" "test" {
  name = "Journal Test Manufacturer"
  slug = "journal-test-manufacturer"
}

resource "netbox_device_role" "test" {
  name  = "Journal Test Role"
  slug  = "journal-test-role"
  color = "ff0000"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Journal Test Model"
  slug         = "journal-test-model"
}

resource "netbox_device" "test" {
  name        = "Journal Test Device"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
}

# Journal entry on a device
resource "netbox_journal_entry" "device_entry" {
  assigned_object_type = "dcim.device"
  assigned_object_id   = netbox_device.test.id
  comments             = "Device provisioned and ready for deployment"
  kind                 = "success"
}

# Outputs for verification
output "basic_id" {
  value = netbox_journal_entry.basic.id
}

output "info_id" {
  value = netbox_journal_entry.info.id
}

output "danger_kind" {
  value = netbox_journal_entry.danger.kind
}

output "device_entry_id" {
  value = netbox_journal_entry.device_entry.id
}
