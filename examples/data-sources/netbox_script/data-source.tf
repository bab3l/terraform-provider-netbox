# Scripts are filesystem-managed in NetBox, so this data source is for reading
# metadata about scripts that already exist on the server.
data "netbox_script" "by_name" {
  name = "sync_devices"
}

data "netbox_script" "by_id" {
  id = "42"
}

output "script_id" {
  value       = data.netbox_script.by_name.id
  description = "The unique ID of the script"
}

output "script_name" {
  value       = data.netbox_script.by_name.name
  description = "The script name"
}

output "script_description" {
  value       = data.netbox_script.by_name.description
  description = "The description shown in NetBox"
}

output "script_is_executable" {
  value       = data.netbox_script.by_name.is_executable
  description = "Whether the script can be executed from NetBox"
}

output "script_module" {
  value       = data.netbox_script.by_id.module
  description = "The module ID containing the script"
}
