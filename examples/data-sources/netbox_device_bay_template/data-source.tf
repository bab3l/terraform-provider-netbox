# Example: Looking up Device Bay Templates from Netbox

# Look up a device bay template by ID
data "netbox_device_bay_template" "by_id" {
  id = 1
}

# Look up a device bay template by name (requires device_type for uniqueness)
data "netbox_device_bay_template" "by_name" {
  name        = "Bay 1"
  device_type = "123"  # Device type ID
}

# Use the data source to reference device bay template properties
output "device_bay_template_details" {
  value = {
    id               = data.netbox_device_bay_template.by_name.id
    name             = data.netbox_device_bay_template.by_name.name
    device_type      = data.netbox_device_bay_template.by_name.device_type
    device_type_name = data.netbox_device_bay_template.by_name.device_type_name
    label            = data.netbox_device_bay_template.by_name.label
    description      = data.netbox_device_bay_template.by_name.description
  }
}
