data "netbox_service" "test" {
  name      = "test-service"
  device_id = 123
}

output "example" {
  value = data.netbox_service.test.id
}
