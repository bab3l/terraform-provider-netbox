data "netbox_device" "test" {
  name = "test-device"
}

output "example" {
  value = data.netbox_device.test.id
}
