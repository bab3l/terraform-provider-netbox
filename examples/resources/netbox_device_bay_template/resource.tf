# Example: Creating Device Bay Templates in Netbox
# Device bay templates define slots for installing child devices in a chassis

# First, create or reference a device type (chassis)
resource "netbox_device_type" "chassis" {
  manufacturer = netbox_manufacturer.example.id
  model        = "Example Blade Chassis"
  slug         = "example-blade-chassis"
  u_height     = 10
  is_full_depth = true
}

# Create a basic device bay template
resource "netbox_device_bay_template" "basic" {
  device_type = netbox_device_type.chassis.id
  name        = "Bay 1"
}

# Create a device bay template with all optional fields
resource "netbox_device_bay_template" "full" {
  device_type = netbox_device_type.chassis.id
  name        = "Bay 2"
  label       = "SLOT-2"
  description = "Second blade slot"
}

# Create multiple device bay templates using count
resource "netbox_device_bay_template" "multi" {
  count       = 8
  device_type = netbox_device_type.chassis.id
  name        = "Blade Slot ${count.index + 1}"
  label       = "BLADE-${count.index + 1}"
  description = "Blade server slot ${count.index + 1}"
}

# Reference an existing manufacturer
resource "netbox_manufacturer" "example" {
  name = "Example Manufacturer"
  slug = "example-manufacturer"
}
