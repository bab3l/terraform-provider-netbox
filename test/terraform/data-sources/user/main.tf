terraform {
  required_version = ">= 1.0"
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {}

# Test data source: Look up user by username
data "netbox_user" "by_username" {
  username = "admin"
}

# Test data source: Look up user by ID (using the ID from the first lookup)
data "netbox_user" "by_id" {
  id = data.netbox_user.by_username.id
}

# Outputs for validation
output "user_id" {
  description = "The ID of the admin user"
  value       = data.netbox_user.by_username.id
}

output "user_username" {
  description = "The username of the admin user"
  value       = data.netbox_user.by_username.username
}

output "id_lookup_username" {
  description = "Username from ID lookup"
  value       = data.netbox_user.by_id.username
}

# Validation outputs
output "all_ids_match" {
  description = "Whether all ID lookups return the same user"
  value       = data.netbox_user.by_username.id == data.netbox_user.by_id.id
}

output "all_usernames_match" {
  description = "Whether all lookups return the same username"
  value       = data.netbox_user.by_username.username == data.netbox_user.by_id.username
}

output "is_admin" {
  description = "Whether the looked up user is admin"
  value       = data.netbox_user.by_username.username == "admin"
}

output "has_valid_id" {
  description = "Whether the user has a valid ID"
  value       = data.netbox_user.by_username.id != ""
}
