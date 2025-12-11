# Look up by ID
data "netbox_custom_link" "by_id" {
  id = "123"
}

# Look up by name
data "netbox_custom_link" "by_name" {
  name = "device_documentation"
}
