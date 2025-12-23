# Look up tunnel termination by ID
data "netbox_tunnel_termination" "by_id" {
  id = "1"
}

# Look up tunnel termination by tunnel ID
data "netbox_tunnel_termination" "by_tunnel" {
  tunnel = "1"
}

# Look up tunnel termination by tunnel name
data "netbox_tunnel_termination" "by_tunnel_name" {
  tunnel_name = "example-tunnel"
}

output "termination_role" {
  value = data.netbox_tunnel_termination.by_id.role
}

output "termination_by_tunnel" {
  value = data.netbox_tunnel_termination.by_tunnel.termination_type
}

output "termination_outside_ip" {
  value = data.netbox_tunnel_termination.by_tunnel_name.outside_ip
}
