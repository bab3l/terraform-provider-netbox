# Create a cable between two interfaces
# Cables represent physical connections in Netbox

# Prerequisites: You need to have devices with interfaces already created
# This example assumes you have two interfaces to connect

resource "netbox_cable" "example" {
  # A-side termination - specify the object type and ID
  a_terminations = [{
    object_type = "dcim.interface"
    object_id   = 1 # ID of interface on first device
  }]

  # B-side termination - specify the object type and ID
  b_terminations = [{
    object_type = "dcim.interface"
    object_id   = 2 # ID of interface on second device
  }]

  # Cable type - common values: cat5e, cat6, cat6a, mmf, smf
  type = "cat6a"

  # Status of the connection
  status = "connected"

  # Physical label on the cable
  label = "PATCH-001"

  # Color code (hex format)
  color = "0000ff"

  # Cable length with unit
  length      = 5.5
  length_unit = "m"

  # Documentation
  description = "Patch cable from switch to server"
  comments    = "Installed during data center buildout"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "cable_vendor"
      value = "Panduit"
    },
    {
      name  = "installation_date"
      value = "2024-01-15"
    },
    {
      name  = "certification_status"
      value = "passed"
    }
  ]

  tags = [
    "patch-cable",
    "certified"
  ]
}

# Example: Fiber cable between interfaces
resource "netbox_cable" "fiber" {
  a_terminations = [{
    object_type = "dcim.interface"
    object_id   = 3
  }]
  b_terminations = [{
    object_type = "dcim.interface"
    object_id   = 4
  }]

  type        = "smf-os2"
  status      = "connected"
  label       = "FIBER-001"
  color       = "ffff00"
  length      = 100
  length_unit = "m"
  description = "Single-mode fiber uplink"

  # Partial custom fields management
  custom_fields = [
    {
      name  = "cable_vendor"
      value = "Corning"
    },
    {
      name  = "fiber_strand_count"
      value = "12"
    },
    {
      name  = "attenuation_db"
      value = "0.35"
    }
  ]

  tags = [
    "fiber-optic",
    "uplink"
  ]
}

# Example: Planned cable (not yet installed)
resource "netbox_cable" "planned" {
  a_terminations = [{
    object_type = "dcim.interface"
    object_id   = 5
  }]
  b_terminations = [{
    object_type = "dcim.interface"
    object_id   = 6
  }]

  status      = "planned"
  type        = "cat6"
  description = "Planned cable for future expansion"
}
