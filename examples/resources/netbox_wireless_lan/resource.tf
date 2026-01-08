resource "netbox_wireless_lan" "test" {
  ssid        = "Test SSID"
  description = "Test wireless LAN for guest access"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "vlan_id"
      value = "100"
    },
    {
      name  = "bandwidth_limit"
      value = "10Mbps"
    },
    {
      name  = "auth_type"
      value = "WPA2-Enterprise"
    },
    {
      name  = "guest_network"
      value = "true"
    }
  ]

  tags = [
    "wireless-lan",
    "guest-network"
  ]
}
