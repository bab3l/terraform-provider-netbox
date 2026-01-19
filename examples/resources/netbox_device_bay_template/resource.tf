# Example: Creating Device Bay Templates in Netbox
# Device bay templates define slots for installing child devices in a chassis

# First, create or reference a device type (chassis)
resource "netbox_device_type" "chassis" {
  manufacturer  = netbox_manufacturer.example.name
  model         = "Example Blade Chassis"
  slug          = "example-blade-chassis"
  u_height      = 10
  is_full_depth = true
}

# Create a basic device bay template
resource "netbox_device_bay_template" "basic" {
  device_type = netbox_device_type.chassis.model
  name        = "Bay 1"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "bay_type"
      value = "compute"
    },
    {
      name  = "max_power_watts"
      value = "150"
    }
  ]
}

# Create a device bay template with all optional fields
resource "netbox_device_bay_template" "full" {
  device_type = netbox_device_type.chassis.model
  name        = "Bay 2"
  label       = "SLOT-2"
  description = "Second blade slot"

  # Partial custom fields management
  custom_fields = [
    {
      name  = "bay_type"
      value = "storage"
    },
    {
      name  = "max_power_watts"
      value = "200"
    }
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_device_bay_template.basic
  id = "123"

  identity = {
    custom_fields = [
      "bay_type:text",
      "max_power_watts:integer",
    ]
  }
}

# Optional: seed owned custom fields during import
import {
  to = netbox_device_bay_template.full
  id = "124"

  identity = {
    custom_fields = [
      "bay_type:text",
      "max_power_watts:integer",
    ]
  }
}

# Create multiple device bay templates using count
resource "netbox_device_bay_template" "multi" {
  count       = 8
  device_type = netbox_device_type.chassis.model
  name        = "Blade Slot ${count.index + 1}"
  label       = "BLADE-${count.index + 1}"
  description = "Blade server slot ${count.index + 1}"
}

# Reference an existing manufacturer
resource "netbox_manufacturer" "example" {
  name = "Example Manufacturer"
  slug = "example-manufacturer"
}
