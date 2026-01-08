resource "netbox_rack_role" "test" {
  name        = "Test Rack Role"
  slug        = "test-rack-role"
  color       = "ff0000"
  description = "Server equipment rack role"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "role_type"
      value = "compute"
    },
    {
      name  = "equipment_category"
      value = "servers"
    },
    {
      name  = "redundancy_level"
      value = "n+1"
    }
  ]

  tags = [
    "rack-role",
    "server-equipment"
  ]
}
