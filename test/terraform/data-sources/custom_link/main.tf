terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

# Create a custom link to look up
resource "netbox_custom_link" "test" {
  name         = "test_link_lookup"
  object_types = ["dcim.device"]
  link_text    = "Test Link"
  link_url     = "https://example.com/{{ object.name }}"
  weight       = 100
  group_name   = "Test Group"
  button_class = "blue"
  new_window   = true
}

# Look up by ID
data "netbox_custom_link" "by_id" {
  id = netbox_custom_link.test.id
}

# Look up by name
data "netbox_custom_link" "by_name" {
  name = netbox_custom_link.test.name
}

output "lookup_by_id_name" {
  value = data.netbox_custom_link.by_id.name
}

output "lookup_by_name_id" {
  value = data.netbox_custom_link.by_name.id
}

output "lookup_by_id_link_text" {
  value = data.netbox_custom_link.by_id.link_text
}

output "lookup_by_name_object_types" {
  value = data.netbox_custom_link.by_name.object_types
}

output "lookup_by_id_button_class" {
  value = data.netbox_custom_link.by_id.button_class
}

output "lookup_by_name_group_name" {
  value = data.netbox_custom_link.by_name.group_name
}
