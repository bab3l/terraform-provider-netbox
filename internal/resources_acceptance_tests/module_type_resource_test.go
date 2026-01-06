package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
