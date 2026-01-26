package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVRFResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vrf")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVRFCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVRFDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVRFResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "name", name),
				),
			},
			{
				Config:   testAccVRFResourceConfig_basic(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccVRFResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vrf-full")
	rd := "65000:100"
	description := "Test VRF with all fields"
	importTargetName := testutil.RandomName("65000:300")
	exportTargetName := testutil.RandomName("65000:400")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVRFCleanup(name)
	cleanup.RegisterRouteTargetCleanup(importTargetName)
	cleanup.RegisterRouteTargetCleanup(exportTargetName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVRFDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVRFResourceConfig_full(name, rd, description, importTargetName, exportTargetName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "name", name),
					resource.TestCheckResourceAttr("netbox_vrf.test", "rd", rd),
					resource.TestCheckResourceAttr("netbox_vrf.test", "description", description),
					resource.TestCheckResourceAttr("netbox_vrf.test", "enforce_unique", "true"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "import_targets.#", "1"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "export_targets.#", "1"),
				),
			},
			{
				Config:   testAccVRFResourceConfig_full(name, rd, description, importTargetName, exportTargetName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccVRFResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vrf-update")
	updatedName := testutil.RandomName("tf-test-vrf-updated")
	importTargetName := testutil.RandomName("65000:310")
	exportTargetName := testutil.RandomName("65000:410")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVRFCleanup(name)
	cleanup.RegisterVRFCleanup(updatedName)
	cleanup.RegisterRouteTargetCleanup(importTargetName)
	cleanup.RegisterRouteTargetCleanup(exportTargetName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVRFDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVRFResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "name", name),
				),
			},
			{
				Config:   testAccVRFResourceConfig_basic(name),
				PlanOnly: true,
			},
			{
				Config: testAccVRFResourceConfig_full(updatedName, "65000:200", "Updated description", importTargetName, exportTargetName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_vrf.test", "rd", "65000:200"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "description", "Updated description"),
				),
			},
			{
				Config:   testAccVRFResourceConfig_full(updatedName, "65000:200", "Updated description", importTargetName, exportTargetName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccVRFResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("tf-test-vrf-extdel")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVRFResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					// Find VRF by name
					items, _, err := client.IpamAPI.IpamVrfsList(context.Background()).Name([]string{name}).Execute()
					if err != nil {
						t.Fatalf("Failed to list VRFs: %v", err)
					}
					if items == nil || len(items.Results) == 0 {
						t.Fatalf("VRF not found with name: %s", name)
					}

					// Delete the VRF
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamVrfsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete VRF: %v", err)
					}

					t.Logf("Successfully externally deleted VRF with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccVRFResource_import(t *testing.T) {
	t.Parallel()

	name := "test-vrf-" + testutil.GenerateSlug("vrf")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVRFResourceConfig_basic(name),
			},
			{
				Config:   testAccVRFResourceConfig_basic(name),
				PlanOnly: true,
			},
			{
				ResourceName:      "netbox_vrf.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccVRFResourceConfig_basic(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccVRFResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vrf-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVRFCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVRFDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVRFResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "name", name),
				),
			},
			{
				Config:   testAccVRFResourceConfig_basic(name),
				PlanOnly: true,
			},
		},
	})
}

func testAccVRFResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_vrf" "test" {
  name = %q
}
`, name)
}

func testAccVRFResourceConfig_full(name, rd, description, importTargetName, exportTargetName string) string {
	return fmt.Sprintf(`
resource "netbox_route_target" "import" {
  name = %q
}

resource "netbox_route_target" "export" {
  name = %q
}

resource "netbox_vrf" "test" {
  name           = %q
  rd             = %q
  description    = %q
  enforce_unique = true
  import_targets = [netbox_route_target.import.id]
  export_targets = [netbox_route_target.export.id]
}
`, importTargetName, exportTargetName, name, rd, description)
}

func testAccVRFResourceConfig_withTargets(name, importTargetName, exportTargetName string) string {
	return fmt.Sprintf(`
resource "netbox_route_target" "import" {
  name = %q
}

resource "netbox_route_target" "export" {
  name = %q
}

resource "netbox_vrf" "test" {
  name           = %q
  import_targets = [netbox_route_target.import.id]
  export_targets = [netbox_route_target.export.id]
}
`, importTargetName, exportTargetName, name)
}

func testAccVRFResourceConfig_withoutTargets(name, importTargetName, exportTargetName string) string {
	return fmt.Sprintf(`
resource "netbox_route_target" "import" {
  name = %q
}

resource "netbox_route_target" "export" {
  name = %q
}

resource "netbox_vrf" "test" {
  name = %q
}
`, importTargetName, exportTargetName, name)
}

func TestAccVRFResource_removeRouteTargets(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vrf-rt-remove")
	importTargetName := testutil.RandomName("65000:320")
	exportTargetName := testutil.RandomName("65000:420")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVRFCleanup(name)
	cleanup.RegisterRouteTargetCleanup(importTargetName)
	cleanup.RegisterRouteTargetCleanup(exportTargetName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVRFDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVRFResourceConfig_withTargets(name, importTargetName, exportTargetName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "import_targets.#", "1"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "export_targets.#", "1"),
				),
			},
			{
				Config: testAccVRFResourceConfig_withoutTargets(name, importTargetName, exportTargetName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),
					resource.TestCheckNoResourceAttr("netbox_vrf.test", "import_targets"),
					resource.TestCheckNoResourceAttr("netbox_vrf.test", "export_targets"),
				),
			},
			{
				Config: testAccVRFResourceConfig_withTargets(name, importTargetName, exportTargetName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "import_targets.#", "1"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "export_targets.#", "1"),
				),
			},
		},
	})
}

func TestAccVRFResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vrf-remove")
	tenantName := testutil.RandomName("tf-test-tenant-remove")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-remove")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVRFCleanup(name)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVRFDestroy,
			testutil.CheckTenantDestroy,
		),
		Steps: []resource.TestStep{
			// Step 1: Create VRF with tenant
			{
				Config: testAccVRFResourceConfig_withTenant(name, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_vrf.test", "tenant"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "description", "VRF with tenant"),
				),
			},
			// Step 2: Remove tenant and verify it's actually removed
			{
				Config: testAccVRFResourceConfig_noTenant(name, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_vrf.test", "tenant"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "description", "VRF after tenant removal"),
				),
			},
			// Step 3: Re-add tenant to verify it can be set again
			{
				Config: testAccVRFResourceConfig_withTenant(name, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_vrf.test", "tenant"),
				),
			},
		},
	})
}

func TestAccConsistency_VRF(t *testing.T) {
	t.Parallel()

	vrfName := testutil.RandomName("vrf")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVRFConsistencyConfig(vrfName, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vrf.test", "name", vrfName),
					testutil.ReferenceFieldCheck("netbox_vrf.test", "tenant"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccVRFConsistencyConfig(vrfName, tenantName, tenantSlug),
			},
		},
	})
}

func testAccVRFConsistencyConfig(vrfName, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = "%[2]s"
  slug = "%[3]s"
}

resource "netbox_vrf" "test" {
  name = "%[1]s"
	tenant = netbox_tenant.test.id
}
`, vrfName, tenantName, tenantSlug)
}

func testAccVRFResourceConfig_withTenant(name, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_vrf" "test" {
  name        = %q
  tenant      = netbox_tenant.test.id
  description = "VRF with tenant"
}
`, tenantName, tenantSlug, name)
}

func testAccVRFResourceConfig_noTenant(name, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_vrf" "test" {
  name        = %q
  description = "VRF after tenant removal"
  # tenant intentionally omitted - should be null in state
}
`, tenantName, tenantSlug, name)
}

func TestAccConsistency_VRF_LiteralNames(t *testing.T) {
	t.Parallel()

	vrfName := testutil.RandomName("vrf-lit")
	tenantName := testutil.RandomName("tenant-lit")
	tenantSlug := testutil.RandomSlug("tenant-lit")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVRFCleanup(vrfName)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVRFDestroy,
			testutil.CheckTenantDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVRFConsistencyLiteralNamesConfig(vrfName, tenantName, tenantSlug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "name", vrfName),
					resource.TestCheckResourceAttr("netbox_vrf.test", "description", description),
				),
			},
			{
				Config:   testAccVRFConsistencyLiteralNamesConfig(vrfName, tenantName, tenantSlug, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),
				),
			},
		},
	})
}

func testAccVRFConsistencyLiteralNamesConfig(vrfName, tenantName, tenantSlug, description string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = "%[2]s"
  slug = "%[3]s"
}

resource "netbox_vrf" "test" {
  name        = "%[1]s"
	tenant      = netbox_tenant.test.id
  description = "%[4]s"
}
`, vrfName, tenantName, tenantSlug, description)
}

func TestAccVRFResource_removeDescriptionAndComments(t *testing.T) {
	t.Parallel()

	vrfName := testutil.RandomName("tf-test-vrf-optional")
	vrfRD := testutil.RandomName("65000:999")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVRFCleanup(vrfRD)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_vrf",
		BaseConfig: func() string {
			return testAccVRFResourceConfig_withRD(vrfName, vrfRD)
		},
		ConfigWithFields: func() string {
			return testAccVRFResourceConfig_withDescriptionAndComments(
				vrfName,
				vrfRD,
				"Test description",
				"Test comments",
			)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
			"comments":    "Test comments",
		},
		RequiredFields: map[string]string{
			"name": vrfName,
		},
		CheckDestroy: testutil.CheckVRFDestroy,
	})
}

func TestAccVRFResource_removeOptionalFields_rd(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vrf-rd-rem")
	rd := "65000:999"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVRFCleanup(name)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_vrf",
		BaseConfig: func() string {
			return testAccVRFResourceConfig_basic(name)
		},
		ConfigWithFields: func() string {
			return testAccVRFResourceConfig_withRD(name, rd)
		},
		OptionalFields: map[string]string{
			"rd": rd,
		},
		RequiredFields: map[string]string{
			"name": name,
		},
		CheckDestroy: testutil.CheckVRFDestroy,
	})
}

func TestAccVRFResource_EnforceUniqueOptionalComputed(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vrf-enforce-unique")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVRFCleanup(name)

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_vrf",
		OptionalField:  "enforce_unique",
		DefaultValue:   "true",
		FieldTestValue: "false",
		BaseConfig: func() string {
			return fmt.Sprintf(`
resource "netbox_vrf" "test" {
  name = %q
}
`, name)
		},
		WithFieldConfig: func(value string) string {
			return fmt.Sprintf(`
resource "netbox_vrf" "test" {
  name           = %q
  enforce_unique = %s
}
`, name, value)
		},
		CheckDestroy: testutil.CheckVRFDestroy,
	})
}

func testAccVRFResourceConfig_withRD(name, rd string) string {
	return fmt.Sprintf(`
resource "netbox_vrf" "test" {
  name = %[1]q
  rd   = %[2]q
}
`, name, rd)
}

func testAccVRFResourceConfig_withDescriptionAndComments(name, rd, description, comments string) string {
	return fmt.Sprintf(`
resource "netbox_vrf" "test" {
  name        = %[1]q
  rd          = %[2]q
  description = %[3]q
  comments    = %[4]q
}
`, name, rd, description, comments)
}

// NOTE: Custom field tests for VRF resource are in resources_acceptance_tests_customfields package

func TestAccVRFResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_vrf",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_vrf" "test" {
  # name missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
