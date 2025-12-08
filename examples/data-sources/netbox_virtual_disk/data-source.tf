# Look up a virtual disk by ID
data "netbox_virtual_disk" "by_id" {
  id = "1"
}

# Look up a virtual disk by name (requires virtual_machine)
data "netbox_virtual_disk" "by_name" {
  name            = "disk0"
  virtual_machine = netbox_virtual_machine.example.id
}

# Use virtual disk data in outputs
output "disk_info" {
  value = {
    id                   = data.netbox_virtual_disk.by_name.id
    name                 = data.netbox_virtual_disk.by_name.name
    size                 = data.netbox_virtual_disk.by_name.size
    virtual_machine_name = data.netbox_virtual_disk.by_name.virtual_machine_name
    description          = data.netbox_virtual_disk.by_name.description
  }
}
