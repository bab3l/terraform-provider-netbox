# Script Data Source Integration Test
# Note: Scripts in NetBox are read-only - they are Python files loaded from the filesystem.
# This test file documents the data source but cannot be run without pre-existing scripts in NetBox.

terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  # Uses NETBOX_SERVER_URL and NETBOX_API_TOKEN environment variables
}

# NOTE: The netbox_script data source requires pre-existing scripts in NetBox.
# Scripts cannot be created via Terraform - they must be Python files in NetBox's scripts directory.
#
# To test this data source manually:
# 1. Create a Python script file in NetBox's SCRIPTS_ROOT directory
# 2. Uncomment and update the data source below with the script's ID or name
#
# Example script lookup by ID:
# data "netbox_script" "example" {
#   id = "1"
# }
#
# Example script lookup by name:
# data "netbox_script" "example" {
#   name = "MyScript"
# }
#
# Available output attributes:
# - id: The unique identifier of the script
# - name: The name of the script
# - module: The module ID containing the script
# - description: Description of the script
# - is_executable: Whether the script is executable
# - display: Display name of the script

# Placeholder output to allow terraform validate/plan to succeed
output "script_data_source_available" {
  value       = true
  description = "The netbox_script data source is available but requires pre-existing scripts in NetBox to test"
}
