# Look up a cable termination by its numeric ID.
# The NetBox API exposes cable terminations as nested objects, so this data source
# is intended for environments where the termination ID is already known.
data "netbox_cable_termination" "by_id" {
  id = "123"
}

output "cable_termination_id" {
  value       = data.netbox_cable_termination.by_id.id
  description = "The unique ID of the cable termination"
}

output "cable_termination_cable" {
  value       = data.netbox_cable_termination.by_id.cable
  description = "The cable this termination belongs to"
}

output "cable_termination_end" {
  value       = data.netbox_cable_termination.by_id.cable_end
  description = "Which end of the cable this termination is attached to"
}

output "cable_termination_object_type" {
  value       = data.netbox_cable_termination.by_id.termination_type
  description = "The NetBox object type on this side of the cable"
}

output "cable_termination_object_id" {
  value       = data.netbox_cable_termination.by_id.termination_id
  description = "The object ID connected by this termination"
}
