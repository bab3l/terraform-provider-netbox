# Tenant Group with Custom Fields Example

terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  server_url = "https://netbox.example.com"
  api_token  = "your-api-token-here"
}

# Create a tenant group with custom fields
resource "netbox_tenant_group" "managed_services" {
  name        = "Managed Services"
  slug        = "managed-services"
  description = "Tenant group for managed service clients"

  tags = [
    {
      name = "managed"
      slug = "managed"
    },
    {
      name = "external"
      slug = "external"
    }
  ]

  custom_fields = [
    {
      name  = "billing_contact"
      type  = "text"
      value = "billing@example.com"
    },
    {
      name  = "service_level"
      type  = "select"
      value = "premium"
    }
  ]
}

# Output the tenant group with custom fields
output "managed_services_info" {
  value = {
    id            = netbox_tenant_group.managed_services.id
    name          = netbox_tenant_group.managed_services.name
    slug          = netbox_tenant_group.managed_services.slug
    description   = netbox_tenant_group.managed_services.description
    custom_fields = netbox_tenant_group.managed_services.custom_fields
    tags          = netbox_tenant_group.managed_services.tags
  }
}
