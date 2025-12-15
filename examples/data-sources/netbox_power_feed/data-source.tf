data "netbox_power_feed" "test" {
  name = "test-power-feed"
}

output "example" {
  value = data.netbox_power_feed.test.id
}
