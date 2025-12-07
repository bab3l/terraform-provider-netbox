# Module Data Source Outputs

output "by_id_device" {
  value = data.netbox_module.by_id.device
}

output "by_id_module_bay" {
  value = data.netbox_module.by_id.module_bay
}

output "by_id_module_type" {
  value = data.netbox_module.by_id.module_type
}

output "by_id_status" {
  value = data.netbox_module.by_id.status
}
