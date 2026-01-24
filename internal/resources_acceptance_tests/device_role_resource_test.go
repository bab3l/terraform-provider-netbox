package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceRoleResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-device-role")
	slug := testutil.RandomSlug("tf-test-dr")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_device_role.test", "slug", slug),
				),
			},
			{
				// Test import
				ResourceName:      "netbox_device_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccDeviceRoleResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccDeviceRoleResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-device-role-full")
	slug := testutil.RandomSlug("tf-test-dr-full")
	description := testutil.RandomName("description")
	configTemplateName := testutil.RandomName("tf-test-config-template")
	configTemplateCode := "{{ device.name }}"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceRoleCleanup(slug)
	cleanup.RegisterConfigTemplateCleanup(configTemplateName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceRoleResourceConfig_full(name, slug, description, "aa1409", false, configTemplateName, configTemplateCode),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_device_role.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_device_role.test", "description", description),
					resource.TestCheckResourceAttr("netbox_device_role.test", "color", "aa1409"),
					resource.TestCheckResourceAttrSet("netbox_device_role.test", "config_template"),
				),
			},
		},
	})
}

func TestAccDeviceRoleResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-device-role-update")
	slug := testutil.RandomSlug("tf-test-dr-upd")
	updatedName := testutil.RandomName("tf-test-device-role-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_role.test", "name", name),
				),
			},
			{
				Config: testAccDeviceRoleResourceConfig_basic(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_role.test", "name", updatedName),
				),
			},
		},
	})
}

func testAccDeviceRoleResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_device_role" "test" {
  name  = %q
  slug  = %q
  color = "aa1409"
}
`, name, slug)
}

func testAccDeviceRoleResourceConfig_baseWithTemplate(name, slug, configTemplateName, configTemplateCode string) string {
	return fmt.Sprintf(`
resource "netbox_config_template" "test" {
	name          = %q
	template_code = %q
}

resource "netbox_device_role" "test" {
	name  = %q
	slug  = %q
	color = "aa1409"
}
`, configTemplateName, configTemplateCode, name, slug)
}

func testAccDeviceRoleResourceConfig_full(name, slug, description, color string, vmRole bool, configTemplateName, configTemplateCode string) string {
	return fmt.Sprintf(`
resource "netbox_config_template" "test" {
	name          = %q
	template_code = %q
}

resource "netbox_device_role" "test" {
  name        = %q
  slug        = %q
  description = %q
  color       = %q
  vm_role     = %t
	config_template = netbox_config_template.test.id
}
`, configTemplateName, configTemplateCode, name, slug, description, color, vmRole)
}

func TestAccConsistency_DeviceRole_LiteralNames(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-device-role-lit")
	slug := testutil.RandomSlug("tf-test-dr-lit")
	description := testutil.RandomName("description")
	color := testutil.Color

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceRoleConsistencyLiteralNamesConfig(name, slug, description, color),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_device_role.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_device_role.test", "description", description),
					resource.TestCheckResourceAttr("netbox_device_role.test", "color", color),
				),
			},
			{
				Config:   testAccDeviceRoleConsistencyLiteralNamesConfig(name, slug, description, color),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_role.test", "id"),
				),
			},
		},
	})
}

func testAccDeviceRoleConsistencyLiteralNamesConfig(name, slug, description, color string) string {
	return fmt.Sprintf(`
resource "netbox_device_role" "test" {
  name        = %q
  slug        = %q
  description = %q
  color       = %q
}
`, name, slug, description, color)
}

func TestAccDeviceRoleResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-device-role-del")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_device_role.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimDeviceRolesList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find device_role for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimDeviceRolesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete device_role: %v", err)
					}
					t.Logf("Successfully externally deleted device_role with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDeviceRoleResource_removeDescription(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-devrole-optional")
	slug := testutil.RandomSlug("tf-test-devrole-optional")
	configTemplateName := testutil.RandomName("tf-test-config-template")
	configTemplateCode := "{{ device_role.name }}"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceRoleCleanup(slug)
	cleanup.RegisterConfigTemplateCleanup(configTemplateName)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_device_role",
		BaseConfig: func() string {
			return testAccDeviceRoleResourceConfig_baseWithTemplate(name, slug, configTemplateName, configTemplateCode)
		},
		ConfigWithFields: func() string {
			return testAccDeviceRoleResourceConfig_withDescription(
				name,
				slug,
				"Test description",
				configTemplateName,
				configTemplateCode,
			)
		},
		OptionalFields: map[string]string{
			"description":     "Test description",
			"config_template": configTemplateName,
		},
		RequiredFields: map[string]string{
			"name": name,
			"slug": slug,
		},
		CheckDestroy: testutil.CheckDeviceRoleDestroy,
	})
}

func testAccDeviceRoleResourceConfig_withDescription(name, slug, description, configTemplateName, configTemplateCode string) string {
	return fmt.Sprintf(`
resource "netbox_config_template" "test" {
	name          = %q
	template_code = %q
}

resource "netbox_device_role" "test" {
	name            = %[3]q
	slug            = %[4]q
	color           = "aa1409"
	description     = %[5]q
	config_template = netbox_config_template.test.name
}
`, configTemplateName, configTemplateCode, name, slug, description)
}

// TestAccDeviceRoleResource_validationErrors tests validation error scenarios.
func TestAccDeviceRoleResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_device_role",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
resource "netbox_device_role" "test" {
  slug = "test-role"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_slug": {
				Config: func() string {
					return `
resource "netbox_device_role" "test" {
  name = "Test Role"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}

// NOTE: Custom field tests for device_role resource are in resources_acceptance_tests_customfields package
