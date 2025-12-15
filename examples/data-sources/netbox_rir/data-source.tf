data "netbox_rir" "test" {
  name = "test-rir"
}

output "example" {
  value = data.netbox_rir.test.id
}
