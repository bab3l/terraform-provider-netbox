data "netbox_user" "admin" {
  username = "admin"
}

output "example" {
  value = data.netbox_user.admin.id
}
