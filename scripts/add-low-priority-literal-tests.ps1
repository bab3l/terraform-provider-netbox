# Script to add LiteralNames consistency tests for low-priority template and device component resources

# FrontPortTemplate - tests device_type and rear_port_template
$frontPortTemplate = @"


// TestAccConsistency_FrontPortTemplate_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_FrontPortTemplate_LiteralNames(t *testing.T) {
	t.Parallel()
	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	rearPortName := testutil.RandomName("rear-port")
	frontPortName := testutil.RandomName("front-port")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFrontPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, frontPortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "name", frontPortName),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "device_type", deviceTypeSlug),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccFrontPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, frontPortName),
			},
		},
	})
}

func testAccFrontPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, frontPortName string) string {
	return fmt.Sprintf(```

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_rear_port_template" "test" {
  device_type = netbox_device_type.test.slug
  name        = %q
  type        = "8p8c"
}

resource "netbox_front_port_template" "test" {
  # Use literal string slugs to mimic existing user state
  device_type       = %q
  rear_port_template = %q
  name              = %q
  type              = "8p8c"

  depends_on = [netbox_device_type.test, netbox_rear_port_template.test]
}

````, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, deviceTypeSlug, rearPortName, frontPortName)
}
"@

# InterfaceTemplate - tests device_type
$interfaceTemplate = @"


// TestAccConsistency_InterfaceTemplate_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_InterfaceTemplate_LiteralNames(t *testing.T) {
	t.Parallel()
	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	interfaceName := testutil.RandomName("interface")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, interfaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_interface_template.test", "name", interfaceName),
					resource.TestCheckResourceAttr("netbox_interface_template.test", "device_type", deviceTypeSlug),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccInterfaceTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, interfaceName),
			},
		},
	})
}

func testAccInterfaceTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, interfaceName string) string {
	return fmt.Sprintf(```

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_interface_template" "test" {
  # Use literal string slug to mimic existing user state
  device_type = %q
  name = %q
  type = "1000base-t"

  depends_on = [netbox_device_type.test]
}

````, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceTypeSlug, interfaceName)
}
"@

# ModuleBayTemplate - tests device_type
$moduleBayTemplate = @"


// TestAccConsistency_ModuleBayTemplate_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_ModuleBayTemplate_LiteralNames(t *testing.T) {
	t.Parallel()
	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	moduleBayName := testutil.RandomName("module-bay")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleBayTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, moduleBayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "name", moduleBayName),
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "device_type", deviceTypeSlug),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccModuleBayTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, moduleBayName),
			},
		},
	})
}

func testAccModuleBayTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, moduleBayName string) string {
	return fmt.Sprintf(```

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_module_bay_template" "test" {
  # Use literal string slug to mimic existing user state
  device_type = %q
  name = %q

  depends_on = [netbox_device_type.test]
}

````, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceTypeSlug, moduleBayName)
}
"@

# PowerOutletTemplate - tests device_type and power_port_template
$powerOutletTemplate = @"


// TestAccConsistency_PowerOutletTemplate_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_PowerOutletTemplate_LiteralNames(t *testing.T) {
	t.Parallel()
	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	powerPortName := testutil.RandomName("power-port")
	outletName := testutil.RandomName("outlet")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerOutletTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, powerPortName, outletName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "name", outletName),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "device_type", deviceTypeSlug),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccPowerOutletTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, powerPortName, outletName),
			},
		},
	})
}

func testAccPowerOutletTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, powerPortName, outletName string) string {
	return fmt.Sprintf(```

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_power_port_template" "test" {
  device_type = netbox_device_type.test.slug
  name        = %q
  type        = "iec-60320-c14"
}

resource "netbox_power_outlet_template" "test" {
  # Use literal string slugs to mimic existing user state
  device_type        = %q
  power_port_template = %q
  name               = %q
  type               = "iec-60320-c13"

  depends_on = [netbox_device_type.test, netbox_power_port_template.test]
}

````, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, powerPortName, deviceTypeSlug, powerPortName, outletName)
}
"@

# PowerPortTemplate - tests device_type
$powerPortTemplate = @"


// TestAccConsistency_PowerPortTemplate_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_PowerPortTemplate_LiteralNames(t *testing.T) {
	t.Parallel()
	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	portName := testutil.RandomName("port")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "name", portName),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "device_type", deviceTypeSlug),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccPowerPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
			},
		},
	})
}

func testAccPowerPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName string) string {
	return fmt.Sprintf(```

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_power_port_template" "test" {
  # Use literal string slug to mimic existing user state
  device_type = %q
  name = %q
  type = "iec-60320-c14"

  depends_on = [netbox_device_type.test]
}

````, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceTypeSlug, portName)
}
"@

# RearPortTemplate - tests device_type
$rearPortTemplate = @"


// TestAccConsistency_RearPortTemplate_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_RearPortTemplate_LiteralNames(t *testing.T) {
	t.Parallel()
	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	portName := testutil.RandomName("port")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRearPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", portName),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "device_type", deviceTypeSlug),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccRearPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
			},
		},
	})
}

func testAccRearPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName string) string {
	return fmt.Sprintf(```

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_rear_port_template" "test" {
  # Use literal string slug to mimic existing user state
  device_type = %q
  name = %q
  type = "8p8c"

  depends_on = [netbox_device_type.test]
}

````, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceTypeSlug, portName)
}
"@

# Apply all template additions
Add-Content -Path "c:\GitRoot\terraform-provider-netbox\internal\resources_test\front_port_template_resource_test.go" -Value $frontPortTemplate
Add-Content -Path "c:\GitRoot\terraform-provider-netbox\internal\resources_test\interface_template_resource_test.go" -Value $interfaceTemplate
Add-Content -Path "c:\GitRoot\terraform-provider-netbox\internal\resources_test\module_bay_template_resource_test.go" -Value $moduleBayTemplate
Add-Content -Path "c:\GitRoot\terraform-provider-netbox\internal\resources_test\power_outlet_template_resource_test.go" -Value $powerOutletTemplate
Add-Content -Path "c:\GitRoot\terraform-provider-netbox\internal\resources_test\power_port_template_resource_test.go" -Value $powerPortTemplate
Add-Content -Path "c:\GitRoot\terraform-provider-netbox\internal\resources_test\rear_port_template_resource_test.go" -Value $rearPortTemplate

Write-Host "Added LiteralNames tests for 6 remaining template resources"
