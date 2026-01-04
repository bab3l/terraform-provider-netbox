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

func TestAccVRFResource_basic(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-vrf")
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVRFCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVRFDestroy,

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
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVRFCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVRFDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccVRFResourceConfig_full(name, rd, description),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "name", name),
					resource.TestCheckResourceAttr("netbox_vrf.test", "rd", rd),
					resource.TestCheckResourceAttr("netbox_vrf.test", "description", description),
					resource.TestCheckResourceAttr("netbox_vrf.test", "enforce_unique", "true"),
				),
			},
			{
				Config:   testAccVRFResourceConfig_full(name, rd, description),
				PlanOnly: true,
			},
		},
	})
}

func TestAccVRFResource_update(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-vrf-update")
	updatedName := testutil.RandomName("tf-test-vrf-updated")
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVRFCleanup(name)
	cleanup.RegisterVRFCleanup(updatedName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVRFDestroy,

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
				Config: testAccVRFResourceConfig_full(updatedName, "65000:200", "Updated description"),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_vrf.test", "rd", "65000:200"),
					resource.TestCheckResourceAttr("netbox_vrf.test", "description", "Updated description"),
				),
			},
			{
				Config:   testAccVRFResourceConfig_full(updatedName, "65000:200", "Updated description"),
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },

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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckVRFDestroy,
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

func testAccVRFResourceConfig_full(name, rd, description string) string {
	return fmt.Sprintf(`
resource "netbox_vrf" "test" {
  name           = %q
  rd             = %q
  description    = %q
  enforce_unique = true
}
`, name, rd, description)
}

func TestAccConsistency_VRF(t *testing.T) {

	t.Parallel()

	vrfName := testutil.RandomName("vrf")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccVRFConsistencyConfig(vrfName, tenantName, tenantSlug),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vrf.test", "name", vrfName),
					resource.TestCheckResourceAttr("netbox_vrf.test", "tenant", tenantName),
				),
			},
			{
				PlanOnly: true,

				Config: testAccVRFConsistencyConfig(vrfName, tenantName, tenantSlug),
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
  tenant = netbox_tenant.test.name
}
`, vrfName, tenantName, tenantSlug)
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
  tenant      = netbox_tenant.test.name
  description = "%[4]s"
}
`, vrfName, tenantName, tenantSlug, description)
}

// NOTE: Custom field tests for VRF resource are in resources_acceptance_tests_customfields package
