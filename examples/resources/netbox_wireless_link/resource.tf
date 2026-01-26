# Basic wireless link between two interfaces
resource "netbox_wireless_link" "example" {
  interface_a = netbox_interface.device_a_wlan0.id
  interface_b = netbox_interface.device_b_wlan0.id
}

# Wireless link with SSID and status
resource "netbox_wireless_link" "with_ssid" {
  interface_a = netbox_interface.device_a_wlan1.id
  interface_b = netbox_interface.device_b_wlan1.id
  ssid        = "MyNetwork"
  status      = "connected"
  description = "Point-to-point wireless link"
}

# Wireless link with authentication
resource "netbox_wireless_link" "with_auth" {
  interface_a = netbox_interface.device_a_wlan2.id
  interface_b = netbox_interface.device_b_wlan2.id
  ssid        = "SecureLink"
  status      = "connected"
  auth_type   = "wpa-personal"
  auth_cipher = "aes"
  auth_psk    = "mysecretpassword"
  description = "Secure wireless link"
}

# Wireless link with distance and all options
resource "netbox_wireless_link" "complete" {
  interface_a   = netbox_interface.device_a_wlan3.id
  interface_b   = netbox_interface.device_b_wlan3.id
  ssid          = "LongRangeLink"
  status        = "connected"
  tenant        = netbox_tenant.example.id
  auth_type     = "wpa-enterprise"
  auth_cipher   = "aes"
  distance      = 5.2
  distance_unit = "km"
  description   = "Long-range point-to-point wireless link"
  comments      = "This link spans across two buildings."

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "frequency_band"
      value = "5GHz"
    },
    {
      name  = "link_capacity"
      value = "1Gbps"
    },
    {
      name  = "signal_strength"
      value = "-65dBm"
    },
    {
      name  = "redundancy_enabled"
      value = "false"
    }
  ]

  tags = [
    "production"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_wireless_link.complete
  id = "123"

  identity = {
    custom_fields = [
      "frequency_band:text",
      "link_capacity:text",
      "signal_strength:text",
      "redundancy_enabled:boolean",
    ]
  }
}
