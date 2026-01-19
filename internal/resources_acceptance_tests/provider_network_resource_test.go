package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProviderNetworkResource_basic(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider")
	providerSlug := testutil.RandomSlug("tf-test-provider")
	networkName := testutil.RandomName("tf-test-network")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderNetworkResourceConfig_basic(providerName, providerSlug, networkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider_network.test", "id"),
					resource.TestCheckResourceAttr("netbox_provider_network.test", "name", networkName),
				),
			},
			{
				ResourceName:      "netbox_provider_network.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccProviderNetworkResourceConfig_basic(providerName, providerSlug, networkName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccProviderNetworkResource_full(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-full")
	providerSlug := testutil.RandomSlug("tf-test-provider-full")
	networkName := testutil.RandomName("tf-test-network-full")
	serviceID := testutil.RandomName("svc")
	description := testutil.RandomName("description")
	updatedDescription := "Updated provider network description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderNetworkResourceConfig_full(providerName, providerSlug, networkName, serviceID, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider_network.test", "id"),
					resource.TestCheckResourceAttr("netbox_provider_network.test", "name", networkName),
					resource.TestCheckResourceAttr("netbox_provider_network.test", "service_id", serviceID),
					resource.TestCheckResourceAttr("netbox_provider_network.test", "description", description),
				),
			},
			{
				Config: testAccProviderNetworkResourceConfig_full(providerName, providerSlug, networkName, serviceID, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider_network.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccProviderNetworkResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-tags")
	providerSlug := testutil.RandomSlug("tf-test-provider-tags")
	networkName := testutil.RandomName("tf-test-network-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderNetworkResourceConfig_tags(providerName, providerSlug, networkName, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider_network.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_provider_network.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_provider_network.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccProviderNetworkResourceConfig_tags(providerName, providerSlug, networkName, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider_network.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_provider_network.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_provider_network.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccProviderNetworkResourceConfig_tags(providerName, providerSlug, networkName, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider_network.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_provider_network.test", "tags.*", tag3Slug),
				),
			},
			{
				Config: testAccProviderNetworkResourceConfig_tags(providerName, providerSlug, networkName, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider_network.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccProviderNetworkResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-tag-order")
	providerSlug := testutil.RandomSlug("tf-test-provider-tag-order")
	networkName := testutil.RandomName("tf-test-network-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderNetworkResourceConfig_tagsOrder(providerName, providerSlug, networkName, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider_network.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_provider_network.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_provider_network.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccProviderNetworkResourceConfig_tagsOrder(providerName, providerSlug, networkName, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider_network.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_provider_network.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_provider_network.test", "tags.*", tag2Slug),
				),
			},
		},
	})
}

func TestAccProviderNetworkResource_update(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-upd")
	providerSlug := testutil.RandomSlug("tf-test-provider-upd")
	networkName := testutil.RandomName("tf-test-network-upd")
	serviceID := "svc-12345"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderNetworkResourceConfig_basic(providerName, providerSlug, networkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider_network.test", "name", networkName),
				),
			},
			{
				Config: testAccProviderNetworkResourceConfig_full(providerName, providerSlug, networkName, serviceID, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider_network.test", "name", networkName),
					resource.TestCheckResourceAttr("netbox_provider_network.test", "service_id", serviceID),
					resource.TestCheckResourceAttr("netbox_provider_network.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccProviderNetworkResource_import(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-imp")
	providerSlug := testutil.RandomSlug("tf-test-provider-imp")
	networkName := testutil.RandomName("tf-test-network-imp")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderNetworkResourceConfig_basic(providerName, providerSlug, networkName),
			},
			{
				ResourceName:      "netbox_provider_network.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             testAccProviderNetworkResourceConfig_basic(providerName, providerSlug, networkName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccProviderNetworkResourceConfig_basic(providerName, providerSlug, networkName string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %q
  slug = %q
}

resource "netbox_provider_network" "test" {
  circuit_provider = netbox_provider.test.id
  name             = %q
}
`, providerName, providerSlug, networkName)
}

func testAccProviderNetworkResourceConfig_tags(providerName, providerSlug, networkName, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag1Uscore2:
		tagsConfig = tagsDoubleSlug
	case caseTag3:
		tagsConfig = tagsSingleSlug
	case tagsEmpty:
		tagsConfig = tagsEmpty
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = "Tag1-%[4]s"
  slug = %[4]q
}

resource "netbox_tag" "tag2" {
  name = "Tag2-%[5]s"
  slug = %[5]q
}

resource "netbox_tag" "tag3" {
  name = "Tag3-%[6]s"
  slug = %[6]q
}

resource "netbox_provider" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_provider_network" "test" {
  circuit_provider = netbox_provider.test.id
  name             = %[3]q
  %[7]s
}
`, providerName, providerSlug, networkName, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccProviderNetworkResourceConfig_tagsOrder(providerName, providerSlug, networkName, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleSlugReversed
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = "Tag1-%[4]s"
  slug = %[4]q
}

resource "netbox_tag" "tag2" {
  name = "Tag2-%[5]s"
  slug = %[5]q
}

resource "netbox_provider" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_provider_network" "test" {
  circuit_provider = netbox_provider.test.id
  name             = %[3]q
  %[6]s
}
`, providerName, providerSlug, networkName, tag1Slug, tag2Slug, tagsConfig)
}

func testAccProviderNetworkResourceConfig_full(providerName, providerSlug, networkName, serviceID, description string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %q
  slug = %q
}

resource "netbox_provider_network" "test" {
  circuit_provider = netbox_provider.test.id
  name             = %q
  service_id       = %q
  description      = %q
}
`, providerName, providerSlug, networkName, serviceID, description)
}

func TestAccConsistency_ProviderNetwork_LiteralNames(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-lit")
	providerSlug := testutil.RandomSlug("tf-test-provider-lit")
	networkName := testutil.RandomName("tf-test-network-lit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderNetworkResourceConfig_basic(providerName, providerSlug, networkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider_network.test", "id"),
					resource.TestCheckResourceAttr("netbox_provider_network.test", "name", networkName),
				),
			},
			{
				Config:   testAccProviderNetworkResourceConfig_basic(providerName, providerSlug, networkName),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider_network.test", "id"),
				),
			},
		},
	})
}

func TestAccProviderNetworkResource_externalDeletion(t *testing.T) {
	t.Parallel()
	providerName := testutil.RandomName("tf-test-provider-ext-del")
	providerSlug := testutil.RandomSlug("provider-ext-del")
	networkName := testutil.RandomName("tf-test-network-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderNetworkResourceConfig_basic(providerName, providerSlug, networkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider_network.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// List provider networks filtered by name
					items, _, err := client.CircuitsAPI.CircuitsProviderNetworksList(context.Background()).NameIc([]string{networkName}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find provider network for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.CircuitsAPI.CircuitsProviderNetworksDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete provider network: %v", err)
					}
					t.Logf("Successfully externally deleted provider network with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccProviderNetworkResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-rem")
	providerSlug := testutil.RandomSlug("tf-test-provider-rem")
	networkName := testutil.RandomName("tf-test-network")
	serviceID := "svc-12345"
	description := "Description"
	comments := "Comments"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderNetworkResourceConfig_fullWithComments(providerName, providerSlug, networkName, serviceID, description, comments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider_network.test", "service_id", serviceID),
					resource.TestCheckResourceAttr("netbox_provider_network.test", "description", description),
					resource.TestCheckResourceAttr("netbox_provider_network.test", "comments", comments),
				),
			},
			{
				Config: testAccProviderNetworkResourceConfig_basic(providerName, providerSlug, networkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("netbox_provider_network.test", "service_id"),
					resource.TestCheckNoResourceAttr("netbox_provider_network.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_provider_network.test", "comments"),
				),
			},
		},
	})
}

func testAccProviderNetworkResourceConfig_fullWithComments(providerName, providerSlug, networkName, serviceID, description, comments string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %q
  slug = %q
}

resource "netbox_provider_network" "test" {
  name             = %q
  circuit_provider = netbox_provider.test.id
  service_id       = %q
  description      = %q
  comments         = %q
}
`, providerName, providerSlug, networkName, serviceID, description, comments)
}

func TestAccProviderNetworkResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_provider_network",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_circuit_provider": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_provider_network" "test" {
  # circuit_provider missing
  name = "Test Network"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_provider" "test" {
  name = "Test Provider"
  slug = "test-provider"
}

resource "netbox_provider_network" "test" {
  circuit_provider = netbox_provider.test.id
  # name missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
