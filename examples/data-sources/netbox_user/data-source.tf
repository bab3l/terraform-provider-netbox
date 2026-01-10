# Look up a user by username
data "netbox_user" "admin" {
  username = "admin"
}

# Look up a user by ID
data "netbox_user" "by_id" {
  id = "1"
}

# Individual attribute outputs
output "user_id" {
  value       = data.netbox_user.admin.id
  description = "The unique ID of the user"
}

output "user_username" {
  value       = data.netbox_user.admin.username
  description = "The username for login"
}

output "user_first_name" {
  value       = data.netbox_user.admin.first_name
  description = "First name of the user"
}

output "user_last_name" {
  value       = data.netbox_user.admin.last_name
  description = "Last name of the user"
}

output "user_email" {
  value       = data.netbox_user.admin.email
  description = "Email address of the user"
}

output "user_is_staff" {
  value       = data.netbox_user.admin.is_staff
  description = "Whether the user has staff privileges"
}

output "user_is_active" {
  value       = data.netbox_user.admin.is_active
  description = "Whether the user account is active"
}

# Access all custom fields
output "user_custom_fields" {
  value       = data.netbox_user.by_id.custom_fields
  description = "All custom fields defined in NetBox for this user"
}

# Access specific custom field by name
output "user_custom_field_example" {
  value       = try([for cf in data.netbox_user.by_id.custom_fields : cf.value if cf.name == "department"][0], null)
  description = "Example: accessing a specific custom field value (department)"
}
