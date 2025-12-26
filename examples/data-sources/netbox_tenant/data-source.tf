# Look up tenant by ID
data "netbox_tenant" "by_id" {
  id = "1"
}

# Look up tenant by slug
data "netbox_tenant" "by_slug" {
  slug = "example-tenant"
}

# Look up tenant by name
data "netbox_tenant" "by_name" {
  name = "Example Tenant"
}

output "tenant_info" {
  value = {
    id          = data.netbox_tenant.by_slug.id
    name        = data.netbox_tenant.by_slug.name
    slug        = data.netbox_tenant.by_slug.slug
    group       = data.netbox_tenant.by_slug.group
    description = data.netbox_tenant.by_slug.description
    comments    = data.netbox_tenant.by_slug.comments
  }
}

output "tenant_by_id" {
  value = data.netbox_tenant.by_id
}

output "tenant_by_name" {
  value = data.netbox_tenant.by_name
}
