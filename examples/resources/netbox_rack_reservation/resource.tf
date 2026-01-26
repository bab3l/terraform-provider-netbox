resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_rack" "test" {
  name   = "Test Rack"
  site   = netbox_site.test.id
  status = "active"
  width  = 19
}

data "netbox_user" "admin" {
  username = "admin"
}

resource "netbox_rack_reservation" "test" {
  rack        = netbox_rack.test.id
  units       = [1, 2, 3]
  user        = data.netbox_user.admin.id
  description = "Reserved for testing"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "project_name"
      value = "Database Migration"
    },
    {
      name  = "reservation_duration_days"
      value = "30"
    },
    {
      name  = "requestor_email"
      value = "dbteam@example.com"
    }
  ]

  tags = [
    "reservation",
    "temporary"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_rack_reservation.test
  id = "123"

  identity = {
    custom_fields = [
      "project_name:text",
      "reservation_duration_days:integer",
      "requestor_email:text",
    ]
  }
}
