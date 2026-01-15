package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceTemplateResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("service-template")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_service_template.test", "protocol", "tcp"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "ports.#", "1"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "ports.0", "80"),
					resource.TestCheckResourceAttrSet("netbox_service_template.test", "id"),
				),
			},
			{
				Config:   testAccServiceTemplateResourceConfig_basic(name),
				PlanOnly: true,
			},
			{
				Config: testAccServiceTemplateResourceConfig_updated(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name+"-updated"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "protocol", "udp"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "ports.#", "2"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "description", "Updated description"),
				),
			},
			{
				Config:   testAccServiceTemplateResourceConfig_updated(name),
				PlanOnly: true,
			},
			{
				ResourceName:      "netbox_service_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccServiceTemplateResourceConfig_updated(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccServiceTemplateResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("service-template")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateResourceConfig_full(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_service_template.test", "protocol", "tcp"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "ports.#", "3"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "description", "Test description"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "comments", "Test comments"),
				),
			},
			{
				Config:   testAccServiceTemplateResourceConfig_full(name),
				PlanOnly: true,
			},
		},
	})
}

func testAccServiceTemplateResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name     = %q
  protocol = "tcp"
  ports    = [80]
}
`, name)
}

func testAccServiceTemplateResourceConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name        = %q
  protocol    = "udp"
  ports       = [53, 123]
  description = "Updated description"
}
`, name+"-updated")
}

func TestAccServiceTemplateResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("service-template-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_service_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
				),
			},
			{
				Config:   testAccServiceTemplateResourceConfig_basic(name),
				PlanOnly: true,
			},
		},
	})
}

func testAccServiceTemplateResourceConfig_full(name string) string {
	return fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name        = %q
  protocol    = "tcp"
  ports       = [80, 443, 8080]
  description = "Test description"
  comments    = "Test comments"
}
`, name)
}

func TestAccServiceTemplateResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("service-template-upd")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateResourceConfig_withDescription(name, testutil.Description1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_service_template.test", "description", testutil.Description1),
				),
			},
			{
				Config:   testAccServiceTemplateResourceConfig_withDescription(name, testutil.Description1),
				PlanOnly: true,
			},
			{
				Config: testAccServiceTemplateResourceConfig_withDescription(name, testutil.Description2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_service_template.test", "description", testutil.Description2),
				),
			},
			{
				Config:   testAccServiceTemplateResourceConfig_withDescription(name, testutil.Description2),
				PlanOnly: true,
			},
		},
	})
}

func testAccServiceTemplateResourceConfig_withDescription(name, description string) string {
	return fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name        = %q
  protocol    = "tcp"
  ports       = [80]
  description = %q
}
`, name, description)
}

func TestAccServiceTemplateResource_external_deletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("service-template-ext")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_service_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
				),
			},
			{
				Config:   testAccServiceTemplateResourceConfig_basic(name),
				PlanOnly: true,
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.IpamAPI.IpamServiceTemplatesList(context.Background()).Name([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find service template for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamServiceTemplatesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete service template: %v", err)
					}
					t.Logf("Successfully externally deleted service template with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
func TestAccConsistency_ServiceTemplate_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("service-template-lit")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateConsistencyLiteralNamesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_service_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
				),
			},
			{
				Config:   testAccServiceTemplateConsistencyLiteralNamesConfig(name),
				PlanOnly: true,
			},
		},
	})
}

func testAccServiceTemplateConsistencyLiteralNamesConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name     = %q
  protocol = "tcp"
  ports    = [80, 443]
}
`, name)
}
func TestAccServiceTemplateResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("service-template-opt")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name        = %[1]q
  protocol    = "udp"
  ports       = [80]
  description = "Description"
  comments    = "Comments"
  tags        = []
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_service_template.test", "protocol", "udp"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "description", "Description"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "comments", "Comments"),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name        = %[1]q
  protocol    = "udp"
  ports       = [80]
  description = "Description"
  comments    = "Comments"
  tags        = []
}
`, name),
				PlanOnly: true,
			},
			{
				Config: fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name     = %[1]q
  ports    = [80]
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_service_template.test", "protocol", "tcp"), // Default protocol is tcp
					resource.TestCheckNoResourceAttr("netbox_service_template.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_service_template.test", "comments"),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name     = %[1]q
  ports    = [80]
}
`, name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccServiceTemplateResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_service_template",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_service_template" "test" {
  # name missing
  ports = [80]
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_ports": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_service_template" "test" {
  name = "test-service"
  # ports missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
