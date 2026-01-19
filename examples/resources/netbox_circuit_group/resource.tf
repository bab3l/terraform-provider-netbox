resource "netbox_circuit_group" "test" {
  name        = "Test Circuit Group"
  slug        = "test-circuit-group"
  description = "Group of related circuits"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "group_type"
      value = "redundant-pair"
    },
    {
      name  = "cost_center"
      value = "IT-NET-001"
    }
  ]

  tags = [
    "circuit-group"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_circuit_group.test
  id = "123"

  identity = {
    custom_fields = [
      "group_type:text",
      "cost_center:text",
    ]
  }
}
