# Example: Basic virtual disk
resource "netbox_virtual_disk" "root" {
  virtual_machine = netbox_virtual_machine.example.name
  name            = "disk0"
  size            = "50" # 50 GB
}

# Example: Virtual disk with description
resource "netbox_virtual_disk" "data" {
  virtual_machine = netbox_virtual_machine.example.name
  name            = "disk1"
  size            = "500"
  description     = "Primary data disk for application storage"
}

# Example: Multiple disks for a VM
resource "netbox_virtual_disk" "os" {
  virtual_machine = netbox_virtual_machine.db.name
  name            = "os-disk"
  size            = "100"
  description     = "Operating system disk"
}

resource "netbox_virtual_disk" "db_data" {
  virtual_machine = netbox_virtual_machine.db.name
  name            = "data-disk"
  size            = "2000"
  description     = "Database data files"
}

resource "netbox_virtual_disk" "db_logs" {
  virtual_machine = netbox_virtual_machine.db.name
  name            = "log-disk"
  size            = "500"
  description     = "Database transaction logs"
}
