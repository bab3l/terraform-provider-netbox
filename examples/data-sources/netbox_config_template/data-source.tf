data "netbox_config_template" "by_id" {
  id = "123"
}

data "netbox_config_template" "by_name" {
  name = "device-config"
}

output "by_id" {
  value = data.netbox_config_template.by_id.name
}

output "by_name" {
  value = data.netbox_config_template.by_name.id
}
