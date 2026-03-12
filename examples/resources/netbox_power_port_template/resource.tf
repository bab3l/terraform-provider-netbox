resource "netbox_manufacturer" "example" {
  name = "Example Manufacturer"
  slug = "example-manufacturer"
}

resource "netbox_device_type" "example" {
  manufacturer = netbox_manufacturer.example.id
  model        = "Example PDU Device Type"
  slug         = "example-pdu-device-type"
}

resource "netbox_power_port_template" "primary_psu" {
  name        = "PSU1"
  device_type = netbox_device_type.example.id
  type        = "iec-60320-c14"
  label       = "Primary PSU Input"
}

resource "netbox_power_port_template" "redundant_psu" {
  name           = "PSU2"
  device_type    = netbox_device_type.example.id
  type           = "iec-60320-c20"
  maximum_draw   = 500
  allocated_draw = 400
  label          = "Redundant PSU Input"
  description    = "Secondary power inlet template for dual-supply devices"
}
