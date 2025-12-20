# Script to add LiteralNames consistency tests for device component resources

# ConsolePort - tests device
$consolePort = @"


// TestAccConsistency_ConsolePort_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_ConsolePort_LiteralNames(t *testing.T) {
	t.Parallel()
	deviceName := testutil.RandomName("device")
	portName := testutil.RandomName("port")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConsolePortConsistencyLiteralNamesConfig(deviceName, portName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_console_port.test", "name", portName),
					resource.TestCheckResourceAttr("netbox_console_port.test", "device", deviceName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccConsolePortConsistencyLiteralNamesConfig(deviceName, portName),
			},
		},
	})
}

func testAccConsolePortConsistencyLiteralNamesConfig(deviceName, portName string) string {
	return fmt.Sprintf(```

resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_role" "test" {
  name = "Test Role"
  slug = "test-role"
}

resource "netbox_device_type" "test" {
  model        = "Test Device Type"
  slug         = "test-device-type"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_console_port" "test" {
  # Use literal string name to mimic existing user state
  device = %q
  name   = %q
  type   = "rj-45"

  depends_on = [netbox_device.test]
}

````, deviceName, deviceName, portName)
}
"@

# ConsoleServerPort - tests device
$consoleServerPort = @"


// TestAccConsistency_ConsoleServerPort_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_ConsoleServerPort_LiteralNames(t *testing.T) {
	t.Parallel()
	deviceName := testutil.RandomName("device")
	portName := testutil.RandomName("port")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConsoleServerPortConsistencyLiteralNamesConfig(deviceName, portName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_console_server_port.test", "name", portName),
					resource.TestCheckResourceAttr("netbox_console_server_port.test", "device", deviceName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccConsoleServerPortConsistencyLiteralNamesConfig(deviceName, portName),
			},
		},
	})
}

func testAccConsoleServerPortConsistencyLiteralNamesConfig(deviceName, portName string) string {
	return fmt.Sprintf(```

resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_role" "test" {
  name = "Test Role"
  slug = "test-role"
}

resource "netbox_device_type" "test" {
  model        = "Test Device Type"
  slug         = "test-device-type"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_console_server_port" "test" {
  # Use literal string name to mimic existing user state
  device = %q
  name   = %q
  type   = "rj-45"

  depends_on = [netbox_device.test]
}

````, deviceName, deviceName, portName)
}
"@

# DeviceBay - tests device
$deviceBay = @"


// TestAccConsistency_DeviceBay_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_DeviceBay_LiteralNames(t *testing.T) {
	t.Parallel()
	deviceName := testutil.RandomName("device")
	bayName := testutil.RandomName("bay")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayConsistencyLiteralNamesConfig(deviceName, bayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device_bay.test", "name", bayName),
					resource.TestCheckResourceAttr("netbox_device_bay.test", "device", deviceName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccDeviceBayConsistencyLiteralNamesConfig(deviceName, bayName),
			},
		},
	})
}

func testAccDeviceBayConsistencyLiteralNamesConfig(deviceName, bayName string) string {
	return fmt.Sprintf(```

resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_role" "test" {
  name = "Test Role"
  slug = "test-role"
}

resource "netbox_device_type" "test" {
  model        = "Test Device Type"
  slug         = "test-device-type"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_device_bay" "test" {
  # Use literal string name to mimic existing user state
  device = %q
  name   = %q

  depends_on = [netbox_device.test]
}

````, deviceName, deviceName, bayName)
}
"@

# FrontPort - tests device and rear_port
$frontPort = @"


// TestAccConsistency_FrontPort_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_FrontPort_LiteralNames(t *testing.T) {
	t.Parallel()
	deviceName := testutil.RandomName("device")
	rearPortName := testutil.RandomName("rear-port")
	frontPortName := testutil.RandomName("front-port")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFrontPortConsistencyLiteralNamesConfig(deviceName, rearPortName, frontPortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_front_port.test", "name", frontPortName),
					resource.TestCheckResourceAttr("netbox_front_port.test", "device", deviceName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccFrontPortConsistencyLiteralNamesConfig(deviceName, rearPortName, frontPortName),
			},
		},
	})
}

func testAccFrontPortConsistencyLiteralNamesConfig(deviceName, rearPortName, frontPortName string) string {
	return fmt.Sprintf(```

resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_role" "test" {
  name = "Test Role"
  slug = "test-role"
}

resource "netbox_device_type" "test" {
  model        = "Test Device Type"
  slug         = "test-device-type"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_rear_port" "test" {
  device = netbox_device.test.name
  name   = %q
  type   = "8p8c"
}

resource "netbox_front_port" "test" {
  # Use literal string names to mimic existing user state
  device    = %q
  rear_port = %q
  name      = %q
  type      = "8p8c"

  depends_on = [netbox_device.test, netbox_rear_port.test]
}

````, deviceName, rearPortName, deviceName, rearPortName, frontPortName)
}
"@

# ModuleBay - tests device
$moduleBay = @"


// TestAccConsistency_ModuleBay_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_ModuleBay_LiteralNames(t *testing.T) {
	t.Parallel()
	deviceName := testutil.RandomName("device")
	moduleBayName := testutil.RandomName("module-bay")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleBayConsistencyLiteralNamesConfig(deviceName, moduleBayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module_bay.test", "name", moduleBayName),
					resource.TestCheckResourceAttr("netbox_module_bay.test", "device", deviceName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccModuleBayConsistencyLiteralNamesConfig(deviceName, moduleBayName),
			},
		},
	})
}

func testAccModuleBayConsistencyLiteralNamesConfig(deviceName, moduleBayName string) string {
	return fmt.Sprintf(```

resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_role" "test" {
  name = "Test Role"
  slug = "test-role"
}

resource "netbox_device_type" "test" {
  model        = "Test Device Type"
  slug         = "test-device-type"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_module_bay" "test" {
  # Use literal string name to mimic existing user state
  device = %q
  name   = %q

  depends_on = [netbox_device.test]
}

````, deviceName, deviceName, moduleBayName)
}
"@

# PowerOutlet - tests device and power_port
$powerOutlet = @"


// TestAccConsistency_PowerOutlet_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_PowerOutlet_LiteralNames(t *testing.T) {
	t.Parallel()
	deviceName := testutil.RandomName("device")
	powerPortName := testutil.RandomName("power-port")
	outletName := testutil.RandomName("outlet")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerOutletConsistencyLiteralNamesConfig(deviceName, powerPortName, outletName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "name", outletName),
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "device", deviceName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccPowerOutletConsistencyLiteralNamesConfig(deviceName, powerPortName, outletName),
			},
		},
	})
}

func testAccPowerOutletConsistencyLiteralNamesConfig(deviceName, powerPortName, outletName string) string {
	return fmt.Sprintf(```

resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_role" "test" {
  name = "Test Role"
  slug = "test-role"
}

resource "netbox_device_type" "test" {
  model        = "Test Device Type"
  slug         = "test-device-type"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_power_port" "test" {
  device = netbox_device.test.name
  name   = %q
  type   = "iec-60320-c14"
}

resource "netbox_power_outlet" "test" {
  # Use literal string names to mimic existing user state
  device     = %q
  power_port = %q
  name       = %q
  type       = "iec-60320-c13"

  depends_on = [netbox_device.test, netbox_power_port.test]
}

````, deviceName, powerPortName, deviceName, powerPortName, outletName)
}
"@

# PowerPort - tests device
$powerPort = @"


// TestAccConsistency_PowerPort_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_PowerPort_LiteralNames(t *testing.T) {
	t.Parallel()
	deviceName := testutil.RandomName("device")
	portName := testutil.RandomName("port")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPortConsistencyLiteralNamesConfig(deviceName, portName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_port.test", "name", portName),
					resource.TestCheckResourceAttr("netbox_power_port.test", "device", deviceName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccPowerPortConsistencyLiteralNamesConfig(deviceName, portName),
			},
		},
	})
}

func testAccPowerPortConsistencyLiteralNamesConfig(deviceName, portName string) string {
	return fmt.Sprintf(```

resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_role" "test" {
  name = "Test Role"
  slug = "test-role"
}

resource "netbox_device_type" "test" {
  model        = "Test Device Type"
  slug         = "test-device-type"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_power_port" "test" {
  # Use literal string name to mimic existing user state
  device = %q
  name   = %q
  type   = "iec-60320-c14"

  depends_on = [netbox_device.test]
}

````, deviceName, deviceName, portName)
}
"@

# RearPort - tests device
$rearPort = @"


// TestAccConsistency_RearPort_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_RearPort_LiteralNames(t *testing.T) {
	t.Parallel()
	deviceName := testutil.RandomName("device")
	portName := testutil.RandomName("port")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRearPortConsistencyLiteralNamesConfig(deviceName, portName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rear_port.test", "name", portName),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "device", deviceName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccRearPortConsistencyLiteralNamesConfig(deviceName, portName),
			},
		},
	})
}

func testAccRearPortConsistencyLiteralNamesConfig(deviceName, portName string) string {
	return fmt.Sprintf(```

resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_role" "test" {
  name = "Test Role"
  slug = "test-role"
}

resource "netbox_device_type" "test" {
  model        = "Test Device Type"
  slug         = "test-device-type"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_rear_port" "test" {
  # Use literal string name to mimic existing user state
  device = %q
  name   = %q
  type   = "8p8c"

  depends_on = [netbox_device.test]
}

````, deviceName, deviceName, portName)
}
"@

# Apply all device component additions
Add-Content -Path "c:\GitRoot\terraform-provider-netbox\internal\resources_test\console_port_resource_test.go" -Value $consolePort
Add-Content -Path "c:\GitRoot\terraform-provider-netbox\internal\resources_test\console_server_port_resource_test.go" -Value $consoleServerPort
Add-Content -Path "c:\GitRoot\terraform-provider-netbox\internal\resources_test\device_bay_resource_test.go" -Value $deviceBay
Add-Content -Path "c:\GitRoot\terraform-provider-netbox\internal\resources_test\front_port_resource_test.go" -Value $frontPort
Add-Content -Path "c:\GitRoot\terraform-provider-netbox\internal\resources_test\module_bay_resource_test.go" -Value $moduleBay
Add-Content -Path "c:\GitRoot\terraform-provider-netbox\internal\resources_test\power_outlet_resource_test.go" -Value $powerOutlet
Add-Content -Path "c:\GitRoot\terraform-provider-netbox\internal\resources_test\power_port_resource_test.go" -Value $powerPort
Add-Content -Path "c:\GitRoot\terraform-provider-netbox\internal\resources_test\rear_port_resource_test.go" -Value $rearPort

Write-Host "Added LiteralNames tests for 8 device component resources"
