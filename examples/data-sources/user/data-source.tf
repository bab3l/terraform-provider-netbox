# Look up a user by username
data "netbox_user" "admin" {
  username = "admin"
}

# Look up a user by ID (you must know the ID beforehand)
data "netbox_user" "by_id" {
  id = "1"
}

# Example: Use the user ID in a rack reservation
resource "netbox_rack_reservation" "example" {
  rack        = netbox_rack.example.id
  units       = [1, 2, 3]
  user        = data.netbox_user.admin.id
  description = "Reserved for ${data.netbox_user.admin.username}"
}
