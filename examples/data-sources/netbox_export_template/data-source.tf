data "netbox_export_template" "test" {
  name = "test-export-template"
}

output "example" {
  value = data.netbox_export_template.test.id
}
