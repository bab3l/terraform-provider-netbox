# Lookup by ID
data "netbox_power_port_template" "by_id" {
  id = 123
}

output "by_id" {
  value = data.netbox_power_port_template.by_id.name
}

# Lookup by device_type and name
data "netbox_power_port_template" "by_device_type" {
  device_type = 456
  name        = "PWR1"
}

output "by_device_type" {
  value = data.netbox_power_port_template.by_device_type.type
}

# Lookup by module_type and name
data "netbox_power_port_template" "by_module_type" {
  module_type = 789
  name        = "PWR1"
}

output "by_module_type" {
  value = data.netbox_power_port_template.by_module_type.type
}
