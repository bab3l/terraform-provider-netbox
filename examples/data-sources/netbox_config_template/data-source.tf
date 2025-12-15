data "netbox_config_template" "test" {
  name = "test-config-template"
}

output "example" {
  value = data.netbox_config_template.test.id
}
