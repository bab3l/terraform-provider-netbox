package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInterfaceTemplateResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-interface-template")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(name + "-mfr-slug")
	cleanup.RegisterDeviceTypeCleanup(name + "-model-slug")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceTemplateResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_interface_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_interface_template.test", "type", "1000base-t"),
				),
			},
			{
				Config:   testAccInterfaceTemplateResourceConfig_basic(name),
				PlanOnly: true,
			},
			{
				// Test import
				ResourceName:            "netbox_interface_template.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device_type", "enabled"},
			},
			{
				Config:   testAccInterfaceTemplateResourceConfig_basic(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccInterfaceTemplateResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-interface-template-full")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(name + "-mfr-slug")
	cleanup.RegisterDeviceTypeCleanup(name + "-model-slug")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceTemplateResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_interface_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_interface_template.test", "type", "1000base-t"),
					resource.TestCheckResourceAttr("netbox_interface_template.test", "enabled", "false"),
					resource.TestCheckResourceAttr("netbox_interface_template.test", "mgmt_only", "true"),
					resource.TestCheckResourceAttr("netbox_interface_template.test", "poe_mode", "pd"),
					resource.TestCheckResourceAttr("netbox_interface_template.test", "poe_type", "type2-ieee802.3at"),
					resource.TestCheckResourceAttr("netbox_interface_template.test", "description", "Test interface template with full options"),
				),
			},
			{
				Config:   testAccInterfaceTemplateResourceConfig_full(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccConsistency_InterfaceTemplate(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-interface-template-consistency")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(name + "-mfr-slug")
	cleanup.RegisterDeviceTypeCleanup(name + "-model-slug")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceTemplateResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_interface_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_interface_template.test", "type", "1000base-t"),
				),
			},
			{
				Config:   testAccInterfaceTemplateResourceConfig_basic(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccConsistency_InterfaceTemplate_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-interface-template-literal")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(name + "-mfr-slug")
	cleanup.RegisterDeviceTypeCleanup(name + "-model-slug")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceTemplateResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_interface_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_interface_template.test", "type", "1000base-t"),
				),
			},
			{
				Config:   testAccInterfaceTemplateResourceConfig_basic(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccInterfaceTemplateResource_external_deletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-interface-template-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(name + "-mfr")
	cleanup.RegisterDeviceTypeCleanup(name + "-dt")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceTemplateResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_interface_template.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimInterfaceTemplatesList(context.Background()).NameIc([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find interface_template for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimInterfaceTemplatesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete interface_template: %v", err)
					}
					t.Logf("Successfully externally deleted interface_template with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccInterfaceTemplateResource_EnabledComprehensive tests comprehensive scenarios for interface template enabled field.
// This validates that Optional+Computed boolean fields work correctly across all scenarios.
func TestAccInterfaceTemplateResource_EnabledComprehensive(t *testing.T) {
	manufacturerName := testutil.RandomName("tf-test-mfr-int-tpl")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-int-tpl")
	deviceTypeName := testutil.RandomName("tf-test-dev-type-int-tpl")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dev-type-int-tpl")
	interfaceTemplateName := testutil.RandomName("tf-test-int-tpl")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_interface_template",
		OptionalField:  "enabled",
		DefaultValue:   "true",
		FieldTestValue: "false",
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckManufacturerDestroy,
		),
		BaseConfig: func() string {
			return testAccInterfaceTemplateResourceWithOptionalField(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, interfaceTemplateName, "enabled", "")
		},
		WithFieldConfig: func(value string) string {
			return testAccInterfaceTemplateResourceWithOptionalField(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, interfaceTemplateName, "enabled", value)
		},
	})
}

func testAccInterfaceTemplateResourceWithOptionalField(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, interfaceTemplateName, optionalFieldName, optionalFieldValue string) string {
	optionalField := ""
	if optionalFieldValue != "" {
		optionalField = fmt.Sprintf("\n  %s = %s", optionalFieldName, optionalFieldValue)
	}

	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_interface_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
  type        = "1000base-t"%s
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, interfaceTemplateName, optionalField)
}

func testAccInterfaceTemplateResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "netbox_interface_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
  type        = "1000base-t"
}
`, testAccInterfaceTemplateResourcePrereqs(name), name)
}

func testAccInterfaceTemplateResourceConfig_full(name string) string {
	return fmt.Sprintf(`
%s

resource "netbox_interface_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
  type        = "1000base-t"
  enabled     = false
  mgmt_only   = true
  poe_mode    = "pd"
  poe_type    = "type2-ieee802.3at"
  description = "Test interface template with full options"
}
`, testAccInterfaceTemplateResourcePrereqs(name), name)
}

func testAccInterfaceTemplateResourcePrereqs(name string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

`, name+"-mfr", name+"-mfr-slug", name+"-model", name+"-model-slug")
}
