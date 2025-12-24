# Look up a VRF by ID
data "netbox_vrf" "by_id" {
  id = "1"
}

# Look up a VRF by name
data "netbox_vrf" "by_name" {
  name = "test-vrf"
}

# Use VRF data in outputs
output "by_id" {
  value = data.netbox_vrf.by_id.name
}

output "by_name" {
  value = data.netbox_vrf.by_name.id
}

output "vrf_info" {
  value = {
    id             = data.netbox_vrf.by_name.id
    name           = data.netbox_vrf.by_name.name
    rd             = data.netbox_vrf.by_name.rd
    tenant         = data.netbox_vrf.by_name.tenant
    enforce_unique = data.netbox_vrf.by_name.enforce_unique
    description    = data.netbox_vrf.by_name.description
  }
}
