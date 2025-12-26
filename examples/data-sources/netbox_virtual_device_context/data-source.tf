# Look up a virtual device context by ID
data "netbox_virtual_device_context" "by_id" {
  id = "1"
}

# Use virtual device context data in outputs
output "vdc_info" {
  value = {
    id          = data.netbox_virtual_device_context.by_id.id
    name        = data.netbox_virtual_device_context.by_id.name
    device      = data.netbox_virtual_device_context.by_id.device
    identifier  = data.netbox_virtual_device_context.by_id.identifier
    status      = data.netbox_virtual_device_context.by_id.status
    primary_ip4 = data.netbox_virtual_device_context.by_id.primary_ip4
    primary_ip6 = data.netbox_virtual_device_context.by_id.primary_ip6
  }
}
