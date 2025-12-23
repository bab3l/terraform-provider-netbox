data "netbox_contact_assignment" "example" {
  id = "123"
}

output "example" {
  value = data.netbox_contact_assignment.example.contact_name
}
