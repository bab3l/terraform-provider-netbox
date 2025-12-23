data "netbox_circuit_type" "by_id" {
  id = "123"
}

data "netbox_circuit_type" "by_name" {
  name = "Internet Transit"
}

data "netbox_circuit_type" "by_slug" {
  slug = "internet-transit"
}

output "by_id" {
  value = data.netbox_circuit_type.by_id.name
}

output "by_name" {
  value = data.netbox_circuit_type.by_name.id
}

output "by_slug" {
  value = data.netbox_circuit_type.by_slug.id
}
