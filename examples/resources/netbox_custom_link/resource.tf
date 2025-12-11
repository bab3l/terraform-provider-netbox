resource "netbox_custom_link" "example" {
  name         = "device_documentation"
  object_types = ["dcim.device"]
  link_text    = "View Documentation"
  link_url     = "https://docs.example.com/devices/{{ object.name }}"
}
