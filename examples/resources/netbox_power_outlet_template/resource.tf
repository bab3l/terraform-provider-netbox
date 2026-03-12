resource "netbox_manufacturer" "example" {
  name = "Example Manufacturer"
  slug = "example-manufacturer"
}

resource "netbox_device_type" "example" {
  manufacturer = netbox_manufacturer.example.id
  model        = "Example PDU Device Type"
  slug         = "example-pdu-device-type"
}

resource "netbox_power_outlet_template" "load_a" {
  name        = "Outlet 1"
  device_type = netbox_device_type.example.id
  type        = "iec-60320-c13"
  label       = "Load A"
}

resource "netbox_power_outlet_template" "load_b" {
  name        = "Outlet 2"
  device_type = netbox_device_type.example.id
  type        = "iec-60320-c19"
  feed_leg    = "B"
  label       = "Load B"
  description = "High-capacity outlet template for redundant feeds"
}
