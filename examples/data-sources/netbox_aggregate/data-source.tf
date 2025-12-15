data "netbox_aggregate" "test" {
  prefix = "10.0.0.0/8"
}

output "example" {
  value = data.netbox_aggregate.test.id
}
