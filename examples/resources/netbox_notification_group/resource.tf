# Create a basic notification group
resource "netbox_notification_group" "ops_team" {
  name        = "ops-team"
  description = "Operations team notification group"
}

# Create a notification group for network engineers
resource "netbox_notification_group" "network_engineers" {
  name        = "network-engineers"
  description = "Network engineering team for infrastructure alerts"
}
