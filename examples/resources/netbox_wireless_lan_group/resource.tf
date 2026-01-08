resource "netbox_wireless_lan_group" "test" {
  name        = "Test WLAN Group"
  slug        = "test-wlan-group"
  description = "Group for organizing wireless LANs by location"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "location_type"
      value = "office"
    },
    {
      name  = "coverage_area"
      value = "Building A"
    },
    {
      name  = "ap_count"
      value = "12"
    }
  ]

  tags = [
    "wlan-group",
    "office"
  ]
}
