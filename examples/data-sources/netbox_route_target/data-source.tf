# Look up a route target by ID
data "netbox_route_target" "by_id" {
  id = "1"
}

# Look up a route target by name
data "netbox_route_target" "by_name" {
  name = "65000:100"
}

# Use route target data in other resources
output "route_target_info" {
  value = {
    id          = data.netbox_route_target.by_name.id
    name        = data.netbox_route_target.by_name.name
    tenant      = data.netbox_route_target.by_name.tenant
    description = data.netbox_route_target.by_name.description
  }
}

output "route_target_by_id" {
  value = data.netbox_route_target.by_id
}
