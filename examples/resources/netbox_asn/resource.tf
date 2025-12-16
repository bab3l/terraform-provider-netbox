resource "netbox_rir" "test" {
  name = "RIPE"
  slug = "ripe"
}

resource "netbox_asn" "test" {
  asn  = 65001
  rir  = netbox_rir.test.name
  tags = []
}
