data "netbox_aggregate" "by_prefix" {
  prefix = "10.0.0.0/8"
}

data "netbox_aggregate" "by_id" {
  id = 123
}

output "by_prefix" {
  value = data.netbox_aggregate.by_prefix.id
}

output "by_id" {
  value = data.netbox_aggregate.by_id.prefix
}
