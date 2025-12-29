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

func TestAccModuleBayTemplateResource_basic(t *testing.T) {

	t.Parallel()
	mfgName := testutil.RandomName("tf-test-mfg")
	mfgSlug := testutil.RandomSlug("tf-test-mfg")
	dtModel := testutil.RandomName("tf-test-dt")
	dtSlug := testutil.RandomSlug("tf-test-dt")
	templateName := testutil.RandomName("tf-test-mbt")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckModuleBayTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleBayTemplateResourceConfig_basic(mfgName, mfgSlug, dtModel, dtSlug, templateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "name", templateName),
				),
			},
			{
				ResourceName:            "netbox_module_bay_template.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device_type"},
			},
		},
	})
}

func TestAccModuleBayTemplateResource_IDPreservation(t *testing.T) {

	t.Parallel()
	mfgName := testutil.RandomName("mfg-mbt")
	mfgSlug := testutil.GenerateSlug(mfgName)
	dtModel := testutil.RandomName("dt-mbt")
	dtSlug := testutil.GenerateSlug(dtModel)
	templateName := testutil.RandomName("mbt-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleBayTemplateResourceConfig_basic(mfgName, mfgSlug, dtModel, dtSlug, templateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "name", templateName),
				),
			},
		},
	})
}

func TestAccModuleBayTemplateResource_full(t *testing.T) {

	t.Parallel()
	mfgName := testutil.RandomName("tf-test-mfg")
	mfgSlug := testutil.RandomSlug("tf-test-mfg")
	dtModel := testutil.RandomName("tf-test-dt")
	dtSlug := testutil.RandomSlug("tf-test-dt")
	templateName := testutil.RandomName("tf-test-mbt")
	label := "Bay 1"
	position := "Front"
	description := "Test module bay template"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckModuleBayTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleBayTemplateResourceConfig_full(mfgName, mfgSlug, dtModel, dtSlug, templateName, label, position, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "name", templateName),
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "label", label),
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "position", position),
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "description", description),
				),
			},
		},
	})
}

func TestAccModuleBayTemplateResource_update(t *testing.T) {

	t.Parallel()
	mfgName := testutil.RandomName("tf-test-mfg")
	mfgSlug := testutil.RandomSlug("tf-test-mfg")
	dtModel := testutil.RandomName("tf-test-dt")
	dtSlug := testutil.RandomSlug("tf-test-dt")
	templateName := testutil.RandomName("tf-test-mbt")
	description1 := testutil.Description1
	description2 := testutil.Description2

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckModuleBayTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleBayTemplateResourceConfig_full(mfgName, mfgSlug, dtModel, dtSlug, templateName, "", "", description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "description", description1),
				),
			},
			{
				Config: testAccModuleBayTemplateResourceConfig_full(mfgName, mfgSlug, dtModel, dtSlug, templateName, "", "", description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "description", description2),
				),
			},
		},
	})
}

func TestAccModuleBayTemplateResource_external_deletion(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("mfg-ext-del")
	mfgSlug := testutil.RandomSlug("mfg-ext-del")
	dtModel := testutil.RandomName("dt-ext-del")
	dtSlug := testutil.RandomSlug("dt-ext-del")
	templateName := testutil.RandomName("mbt-ext-del")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleBayTemplateResourceConfig_basic(mfgName, mfgSlug, dtModel, dtSlug, templateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "name", templateName),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimModuleBayTemplatesList(context.Background()).NameIc([]string{templateName}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find module_bay_template for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimModuleBayTemplatesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete module_bay_template: %v", err)
					}
					t.Logf("Successfully externally deleted module_bay_template with ID: %d", itemID)
				},
				Config: testAccModuleBayTemplateResourceConfig_basic(mfgName, mfgSlug, dtModel, dtSlug, templateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_bay_template.test", "id"),
				),
			},
		},
	})
}

func testAccModuleBayTemplateResourceConfig_basic(mfgName, mfgSlug, dtModel, dtSlug, templateName string) string {
	return fmt.Sprintf(`
provider "netbox" {}

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
  name        = %q
  device_type = netbox_device_type.test.id
}
`, mfgName, mfgSlug, dtModel, dtSlug, templateName)
}

func testAccModuleBayTemplateResourceConfig_full(mfgName, mfgSlug, dtModel, dtSlug, templateName, label, position, description string) string {
	labelAttr := ""
	if label != "" {
		labelAttr = fmt.Sprintf(`label       = %q`, label)
	}
	positionAttr := ""
	if position != "" {
		positionAttr = fmt.Sprintf(`position    = %q`, position)
	}

	descAttr := ""
	if description != "" {
		descAttr = fmt.Sprintf(`description = %q`, description)
	}
	return fmt.Sprintf(`
provider "netbox" {}

resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  model        = %[3]q
  slug         = %[4]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_module_bay_template" "test" {
  name        = %[5]q
  device_type = netbox_device_type.test.id
  %[6]s
  %[7]s
  %[8]s
}
`, mfgName, mfgSlug, dtModel, dtSlug, templateName, labelAttr, positionAttr, descAttr)
}

func TestAccConsistency_ModuleBayTemplate_LiteralNames(t *testing.T) {
	t.Parallel()
	mfgName := testutil.RandomName("tf-test-mfg-lit")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-lit")
	dtModel := testutil.RandomName("tf-test-dt-lit")
	dtSlug := testutil.RandomSlug("tf-test-dt-lit")
	templateName := testutil.RandomName("tf-test-mbt-lit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckModuleBayTemplateDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccModuleBayTemplateConsistencyLiteralNamesConfig(mfgName, mfgSlug, dtModel, dtSlug, templateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "name", templateName),
				),
			},
			{
				Config:   testAccModuleBayTemplateConsistencyLiteralNamesConfig(mfgName, mfgSlug, dtModel, dtSlug, templateName),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_bay_template.test", "id"),
				),
			},
		},
	})
}

func testAccModuleBayTemplateConsistencyLiteralNamesConfig(mfgName, mfgSlug, dtModel, dtSlug, templateName string) string {
	return fmt.Sprintf(`
provider "netbox" {}

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
  name        = %q
  device_type = netbox_device_type.test.id
}
`, mfgName, mfgSlug, dtModel, dtSlug, templateName)
}
