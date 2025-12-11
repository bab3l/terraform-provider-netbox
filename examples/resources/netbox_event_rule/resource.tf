# Create a webhook for the event rule to trigger
resource "netbox_webhook" "device_changes" {
  name        = "device-changes-webhook"
  payload_url = "https://my-automation-server.example.com/webhook"
}

# Create an event rule that triggers on device changes
resource "netbox_event_rule" "device_changes" {
  name               = "device-changes"
  description        = "Trigger webhook when devices are created, updated, or deleted"
  object_types       = ["dcim.device"]
  event_types        = ["object_created", "object_updated", "object_deleted"]
  action_type        = "webhook"
  action_object_type = "extras.webhook"
  action_object_id   = netbox_webhook.device_changes.id
  enabled            = true
}

# Event rule for IP address changes with multiple object types
resource "netbox_event_rule" "ip_changes" {
  name               = "ip-address-changes"
  description        = "Monitor IP address and prefix changes"
  object_types       = ["ipam.ipaddress", "ipam.prefix"]
  event_types        = ["object_created", "object_updated", "object_deleted"]
  action_type        = "webhook"
  action_object_type = "extras.webhook"
  action_object_id   = netbox_webhook.device_changes.id
  enabled            = true
}
