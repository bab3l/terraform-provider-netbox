# Look up service by ID
data "netbox_service" "by_id" {
  id = "1"
}

# Look up service by name and device
data "netbox_service" "by_name_and_device" {
  name   = "http"
  device = "1"
}

# Look up service by name and virtual machine
data "netbox_service" "by_name_and_vm" {
  name            = "ssh"
  virtual_machine = "1"
}

output "service_by_id" {
  value = data.netbox_service.by_id.name
}

output "service_protocol" {
  value = data.netbox_service.by_name_and_device.protocol
}

output "service_ports" {
  value = data.netbox_service.by_name_and_device.ports
}
