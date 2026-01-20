resource "netbox_virtual_chassis" "test" {
  name        = "Test Virtual Chassis"
  domain      = "test-domain"
  description = "Stacked switch chassis"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "stack_protocol"
      value = "VSS"
    },
    {
      name  = "stack_priority"
      value = "150"
    },
    {
      name  = "member_count"
      value = "4"
    }
  ]

  tags = [
    "virtual-chassis",
    "switch-stack"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_virtual_chassis.test
  id = "123"

  identity = {
    custom_fields = [
      "stack_protocol:text",
      "stack_priority:integer",
      "member_count:integer",
    ]
  }
}
