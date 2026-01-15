package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccModuleTypeResource_basic(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg")
	mfgSlug := testutil.RandomSlug("tf-test-mfg")
	model := testutil.RandomName("tf-test-module-type")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleTypeResourceConfig_basic(mfgName, mfgSlug, model),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_module_type.test", "model", model),
				),
			},
			{
				ResourceName:            "netbox_module_type.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"manufacturer"},
			},
		},
	})
}

func TestAccModuleTypeResource_full(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg-full")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-full")
	model := testutil.RandomName("tf-test-module-type-full")
	description := testutil.RandomName("description")
	updatedDescription := "Updated module type description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleTypeResourceConfig_full(mfgName, mfgSlug, model, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_module_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_module_type.test", "part_number", "MT-001"),
					resource.TestCheckResourceAttr("netbox_module_type.test", "description", description),
				),
			},
			{
				Config: testAccModuleTypeResourceConfig_full(mfgName, mfgSlug, model, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module_type.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccModuleTypeResource_IDPreservation(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg-id")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-id")
	model := testutil.RandomName("tf-test-module-type-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleTypeResourceConfig_basic(mfgName, mfgSlug, model),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_module_type.test", "model", model),
					resource.TestCheckResourceAttrSet("netbox_module_type.test", "manufacturer"),
				),
			},
		},
	})
}

func testAccModuleTypeResourceConfig_basic(mfgName, mfgSlug, model string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
}
`, mfgName, mfgSlug, model)
}

func testAccModuleTypeResourceConfig_full(mfgName, mfgSlug, model, description string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  part_number  = "MT-001"
  description  = %q
}
`, mfgName, mfgSlug, model, description)
}

func TestAccConsistency_ModuleType_LiteralNames(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg-lit")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-lit")
	model := testutil.RandomName("tf-test-module-type-lit")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleTypeConsistencyLiteralNamesConfig(mfgName, mfgSlug, model, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_module_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_module_type.test", "description", description),
				),
			},
			{
				Config:   testAccModuleTypeConsistencyLiteralNamesConfig(mfgName, mfgSlug, model, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_type.test", "id"),
				),
			},
		},
	})
}

func testAccModuleTypeConsistencyLiteralNamesConfig(mfgName, mfgSlug, model, description string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  description  = %q
}
`, mfgName, mfgSlug, model, description)
}

func TestAccModuleTypeResource_update(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg-update")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-update")
	model := testutil.RandomName("tf-test-module-type-update")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleTypeResourceConfig_update(mfgName, mfgSlug, model, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_module_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_module_type.test", "comments", testutil.Description1),
				),
			},
			{
				Config: testAccModuleTypeResourceConfig_update(mfgName, mfgSlug, model, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module_type.test", "comments", testutil.Description2),
				),
			},
		},
	})
}

func testAccModuleTypeResourceConfig_update(mfgName, mfgSlug, model, comments string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  comments     = %q
}
`, mfgName, mfgSlug, model, comments)
}

func TestAccModuleTypeResource_external_deletion(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg-ext-del")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-ext-del")
	model := testutil.RandomName("tf-test-module-type-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleTypeResourceConfig_basic(mfgName, mfgSlug, model),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_module_type.test", "model", model),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					ctx := context.Background()

					// Find the module type by model name
					listResp, _, err := client.DcimAPI.DcimModuleTypesList(ctx).Model([]string{model}).Execute()
					if err != nil {
						t.Fatalf("Failed to list module types: %v", err)
					}

					if listResp.Count == 0 {
						t.Fatalf("Module type with model %q not found", model)
					}

					moduleTypeID := listResp.Results[0].Id

					// Delete the module type via API
					_, err = client.DcimAPI.DcimModuleTypesDestroy(ctx, moduleTypeID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete module type: %v", err)
					}

					t.Logf("Successfully externally deleted module type with ID: %d", moduleTypeID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccModuleTypeResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg-opt")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-opt")
	model := testutil.RandomName("tf-test-mt-opt")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[3]q
  description  = "Test module type"
  part_number  = "MT-98765"
  airflow      = "front-to-rear"
  weight       = 2.5
  weight_unit  = "lb"
}
`, mfgName, mfgSlug, model),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module_type.test", "description", "Test module type"),
					resource.TestCheckResourceAttr("netbox_module_type.test", "part_number", "MT-98765"),
					resource.TestCheckResourceAttr("netbox_module_type.test", "airflow", "front-to-rear"),
					resource.TestCheckResourceAttr("netbox_module_type.test", "weight", "2.5"),
					resource.TestCheckResourceAttr("netbox_module_type.test", "weight_unit", "lb"),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[3]q
}
`, mfgName, mfgSlug, model),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("netbox_module_type.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_module_type.test", "part_number"),
					resource.TestCheckNoResourceAttr("netbox_module_type.test", "airflow"),
					resource.TestCheckNoResourceAttr("netbox_module_type.test", "weight"),
					resource.TestCheckNoResourceAttr("netbox_module_type.test", "weight_unit"),
				),
			},
		},
	})
}

func TestAccModuleTypeResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_module_type",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_manufacturer": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_module_type" "test" {
  # manufacturer missing
  model = "test-module-model"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_model": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_manufacturer" "test" {
  name = "test-manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  # model missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
