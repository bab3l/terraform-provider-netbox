# Look up tunnel group by ID
data "netbox_tunnel_group" "by_id" {
  id = "1"
}

# Look up tunnel group by slug
data "netbox_tunnel_group" "by_slug" {
  slug = "default"
}

# Look up tunnel group by name
data "netbox_tunnel_group" "by_name" {
  name = "Default"
}

output "tunnel_group_description" {
  value = data.netbox_tunnel_group.by_id.description
}

output "tunnel_group_by_slug" {
  value = data.netbox_tunnel_group.by_slug.name
}

output "tunnel_group_by_name" {
  value = data.netbox_tunnel_group.by_name.slug
}
