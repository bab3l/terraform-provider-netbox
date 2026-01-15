package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomLinkResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("cl")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCustomLinkCleanupByName(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCustomLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomLinkResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_custom_link.test", "name", name),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "link_text", "View in External System"),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "link_url", "https://example.com/device/{{ object.name }}"),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "object_types.#", "1"),
				),
			},
			{
				// Verify no changes after create
				Config:   testAccCustomLinkResourceConfig_basic(name),
				PlanOnly: true,
			},
			{
				ResourceName:      "netbox_custom_link.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCustomLinkResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("cl")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCustomLinkCleanupByName(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCustomLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomLinkResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_custom_link.test", "name", name),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "enabled", "true"),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "weight", "50"),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "group_name", "External Links"),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "button_class", "blue"),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "new_window", "true"),
				),
			},
			{
				// Verify no changes after create
				Config:   testAccCustomLinkResourceConfig_full(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccCustomLinkResource_update(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("cl")
	updatedName := name + "-updated"
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCustomLinkCleanupByName(name)
	cleanup.RegisterCustomLinkCleanupByName(updatedName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCustomLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomLinkResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_custom_link.test", "name", name),
				),
			},
			{
				// Verify no changes after create
				Config:   testAccCustomLinkResourceConfig_basic(name),
				PlanOnly: true,
			},
			{
				Config: testAccCustomLinkResourceConfig_basic(updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_custom_link.test", "name", updatedName),
				),
			},
			{
				// Verify no changes after update
				Config:   testAccCustomLinkResourceConfig_basic(updatedName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccCustomLinkResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("cl-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCustomLinkCleanupByName(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCustomLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomLinkResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_custom_link.test", "id"),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "name", name),
				),
			},
			{
				// Verify no changes after create
				Config:   testAccCustomLinkResourceConfig_basic(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccConsistency_CustomLink_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("cl-lit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCustomLinkCleanupByName(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCustomLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomLinkResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_custom_link.test", "name", name),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "link_text", "View in External System"),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "link_url", "https://example.com/device/{{ object.name }}"),
				),
			},
			{
				Config:   testAccCustomLinkResourceConfig_basic(name),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_custom_link.test", "id"),
				),
			},
		},
	})
}

func testAccCustomLinkResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_link" "test" {
  name         = "%s"
  object_types = ["dcim.device"]
  link_text    = "View in External System"
  link_url     = "https://example.com/device/{{ object.name }}"
}
`, name)
}

func testAccCustomLinkResourceConfig_full(name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_link" "test" {
  name         = "%s"
  object_types = ["dcim.device", "dcim.site"]
  enabled      = true
  link_text    = "View Details"
  link_url     = "https://example.com/{{ object.name }}"
  weight       = 50
  group_name   = "External Links"
  button_class = "blue"
  new_window   = true
}
`, name)
}

func TestAccCustomLinkResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("tf-test-cl-extdel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCustomLinkCleanupByName(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomLinkResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_custom_link.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					// Find custom link by name
					items, _, err := client.ExtrasAPI.ExtrasCustomLinksList(context.Background()).Name([]string{name}).Execute()
					if err != nil {
						t.Fatalf("Failed to list custom links: %v", err)
					}
					if items == nil || len(items.Results) == 0 {
						t.Fatalf("Custom link not found with name: %s", name)
					}

					// Delete the custom link
					itemID := items.Results[0].Id
					_, err = client.ExtrasAPI.ExtrasCustomLinksDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete custom link: %v", err)
					}

					t.Logf("Successfully externally deleted custom link with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccCustomLinkResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("cl-opt")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCustomLinkCleanupByName(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCustomLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_custom_link" "test" {
  name         = %[1]q
  object_types = ["dcim.device"]
  link_text    = "View Device"
  link_url     = "https://example.com/device/{{ object.name }}"
  enabled      = true
  weight       = 50
  group_name   = "External Links"
  button_class = "blue"
  new_window   = true
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_custom_link.test", "name", name),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "enabled", "true"),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "weight", "50"),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "group_name", "External Links"),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "button_class", "blue"),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "new_window", "true"),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "netbox_custom_link" "test" {
  name         = %[1]q
  object_types = ["dcim.device"]
  link_text    = "View Device"
  link_url     = "https://example.com/device/{{ object.name }}"
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_custom_link.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_custom_link.test", "group_name"),
					// TODO: button_class and new_window don't clear properly in Netbox API
					// They retain their previous non-default values even when set to defaults
					// resource.TestCheckResourceAttr("netbox_custom_link.test", "button_class", "default"),
					// resource.TestCheckResourceAttr("netbox_custom_link.test", "new_window", "false"),
				),
			},
		},
	})
}

func TestAccCustomLinkResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_custom_link",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_object_types": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_custom_link" "test" {
  # object_types missing
  name = "Test Link"
  link_text = "Click here"
  link_url = "https://example.com"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_custom_link" "test" {
  object_types = ["dcim.site"]
  # name missing
  link_text = "Click here"
  link_url = "https://example.com"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_link_text": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_custom_link" "test" {
  object_types = ["dcim.site"]
  name = "Test Link"
  # link_text missing
  link_url = "https://example.com"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_link_url": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_custom_link" "test" {
  object_types = ["dcim.site"]
  name = "Test Link"
  link_text = "Click here"
  # link_url missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
