# Look up site group by ID
data "netbox_site_group" "by_id" {
  id = "1"
}

# Look up site group by slug
data "netbox_site_group" "by_slug" {
  slug = "north-america"
}

# Look up site group by name
data "netbox_site_group" "by_name" {
  name = "North America"
}

output "site_group_parent" {
  value = data.netbox_site_group.by_id.parent
}

output "site_group_description" {
  value = data.netbox_site_group.by_slug.description
}

output "site_group_name" {
  value = data.netbox_site_group.by_name.name
}
