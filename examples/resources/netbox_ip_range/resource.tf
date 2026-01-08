resource "netbox_ip_range" "test" {
  start_address = "10.0.0.1/24"
  end_address   = "10.0.0.10/24"
  status        = "active"
  description   = "DHCP pool for guest WiFi"
  comments      = "Reserved IP range for dynamic allocation"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "range_purpose"
      value = "dhcp-pool"
    },
    {
      name  = "allocation_type"
      value = "dynamic"
    },
    {
      name  = "lease_time_hours"
      value = "24"
    }
  ]

  tags = [
    "dhcp",
    "guest-network"
  ]
}
