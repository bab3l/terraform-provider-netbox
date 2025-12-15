data "netbox_power_panel" "test" {
  name = "test-power-panel"
}

output "example" {
  value = data.netbox_power_panel.test.id
}
