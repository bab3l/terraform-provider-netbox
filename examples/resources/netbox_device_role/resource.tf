resource "netbox_device_role" "test" {
  name        = "Test Role"
  slug        = "test-role"
  color       = "ff0000"
  vm_role     = false
  description = "Network switch role"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "role_category"
      value = "network"
    },
    {
      name  = "support_tier"
      value = "tier-1"
    },
    {
      name  = "monitoring_profile"
      value = "network-devices"
    }
  ]

  tags = [
    "network-role",
    "production"
  ]
}
