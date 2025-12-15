resource "netbox_rir" "test" {
  name = "RIPE"
  slug = "ripe"
}

resource "netbox_asn" "test" {
  asn    = 65001
  rir_id = netbox_rir.test.id
  tags   = []
}
