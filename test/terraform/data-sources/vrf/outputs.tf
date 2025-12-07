# VRF Data Source Outputs

output "by_id_name" {
  value = data.netbox_vrf.by_id.name
}

output "by_id_rd" {
  value = data.netbox_vrf.by_id.rd
}

output "by_name_id" {
  value = data.netbox_vrf.by_name.id
}

output "by_name_description" {
  value = data.netbox_vrf.by_name.description
}
