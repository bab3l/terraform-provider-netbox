#!/bin/bash

# Import an existing tenant group by ID
terraform import netbox_tenant_group.example 123

# Note: Replace "123" with the actual ID of the tenant group you want to import
# You can find the ID in Netbox's web interface or API

# After importing, you'll need to write a configuration that matches the imported resource
# Example:
#
# resource "netbox_tenant_group" "example" {
#   name        = "Example Tenant Group"
#   slug        = "example-tenant-group"
#   description = "An imported tenant group"
# }
