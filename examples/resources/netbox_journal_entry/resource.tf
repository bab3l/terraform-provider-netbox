# Create a journal entry for a site
resource "netbox_journal_entry" "site_update" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.example.id
  kind                 = "info"
  comments             = "Site configuration updated for Q4 deployment."
}

# Journal entry with warning status
resource "netbox_journal_entry" "maintenance_notice" {
  assigned_object_type = "dcim.device"
  assigned_object_id   = netbox_device.example.id
  kind                 = "warning"
  comments             = "Scheduled maintenance window: 2024-01-15 02:00-04:00 UTC"
}

# Optional: seed owned custom fields during import
import {
  to = netbox_journal_entry.site_update
  id = "123"

  identity = {
    custom_fields = [
      "ticket_id:text",
      "change_window:text",
    ]
  }
}
