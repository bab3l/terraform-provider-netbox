resource "netbox_rir" "test" {
  name = "Test RIR"
  slug = "test-rir"
}

resource "netbox_aggregate" "test" {
  prefix      = "10.0.0.0/8"
  rir         = netbox_rir.test.id
  description = "Private IP space allocation"
  comments    = "RFC1918 private address space"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "allocation_type"
      value = "internal-use"
    },
    {
      name  = "allocation_date"
      value = "2024-01-15"
    },
    {
      name  = "utilization_threshold"
      value = "80"
    }
  ]

  tags = [
    "private-space",
    "rfc1918"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_aggregate.test
  id = "123"

  identity = {
    custom_fields = [
      "allocation_type:text",
      "allocation_date:date",
      "utilization_threshold:integer",
    ]
  }
}
