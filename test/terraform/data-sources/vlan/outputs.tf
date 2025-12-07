# VLAN Data Source Outputs

output "by_id_vid" {
  value = data.netbox_vlan.by_id.vid
}

output "by_id_name" {
  value = data.netbox_vlan.by_id.name
}

output "by_id_status" {
  value = data.netbox_vlan.by_id.status
}

output "by_id_description" {
  value = data.netbox_vlan.by_id.description
}
