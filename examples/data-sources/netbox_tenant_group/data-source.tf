# Look up tenant group by ID
data "netbox_tenant_group" "by_id" {
  id = "1"
}

# Look up tenant group by slug
data "netbox_tenant_group" "by_slug" {
  slug = "example-group"
}

# Look up tenant group by name
data "netbox_tenant_group" "by_name" {
  name = "Example Group"
}

output "tenant_group_parent" {
  value = data.netbox_tenant_group.by_id.parent
}

output "tenant_group_description" {
  value = data.netbox_tenant_group.by_slug.description
}

output "tenant_group_name" {
  value = data.netbox_tenant_group.by_name.name
}
