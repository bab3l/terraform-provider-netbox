data "netbox_platform" "test" {
  name = "test-platform"
}

output "example" {
  value = data.netbox_platform.test.id
}
