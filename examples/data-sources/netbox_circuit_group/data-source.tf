data "netbox_circuit_group" "by_id" {
  id = "123"
}

data "netbox_circuit_group" "by_name" {
  name = "example-group"
}

data "netbox_circuit_group" "by_slug" {
  slug = "example-group"
}

output "by_id" {
  value = data.netbox_circuit_group.by_id.name
}

output "by_name" {
  value = data.netbox_circuit_group.by_name.id
}

output "by_slug" {
  value = data.netbox_circuit_group.by_slug.id
}
