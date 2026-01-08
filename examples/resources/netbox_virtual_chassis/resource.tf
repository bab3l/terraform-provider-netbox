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
