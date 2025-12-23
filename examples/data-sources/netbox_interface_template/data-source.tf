data "netbox_interface_template" "by_id" {
  id = 1
}

data "netbox_interface_template" "by_name" {
  name        = "GigabitEthernet"
  device_type = 5
}

data "netbox_interface_template" "by_name_module" {
  name        = "SFP"
  module_type = 10
}

output "template_id" {
  value = data.netbox_interface_template.by_id.id
}

output "template_name" {
  value = data.netbox_interface_template.by_name.name
}

output "template_type" {
  value = data.netbox_interface_template.by_name.type
}
