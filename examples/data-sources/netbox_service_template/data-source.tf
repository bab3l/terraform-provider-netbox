data "netbox_service_template" "test" {
  name = "test-service-template"
}

output "example" {
  value = data.netbox_service_template.test.id
}
