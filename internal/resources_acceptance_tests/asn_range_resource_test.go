package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Acceptance Tests.

func TestAccASNRangeResource_basic(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	name := testutil.RandomName("tf-test-asn-range")

	slug := testutil.RandomSlug("tf-test-asn-range")

	rirName := testutil.RandomName("tf-test-rir")

	rirSlug := testutil.RandomSlug("tf-test-rir")

	// Register cleanup to ensure resources are deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterASNRangeCleanup(name)

	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckASNRangeDestroy,

			testutil.CheckRIRDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "64512"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "64612"),

					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "rir"),
				),
			},
		},
	})

}

func TestAccASNRangeResource_full(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-asn-range-full")

	slug := testutil.RandomSlug("tf-test-asn-range-full")

	rirName := testutil.RandomName("tf-test-rir")

	rirSlug := testutil.RandomSlug("tf-test-rir")

	tenantName := testutil.RandomName("tf-test-tenant")

	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterASNRangeCleanup(name)

	cleanup.RegisterRIRCleanup(rirSlug)

	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckASNRangeDestroy,

			testutil.CheckRIRDestroy,

			testutil.CheckTenantDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccASNRangeResourceConfig_full(name, slug, rirName, rirSlug, tenantName, tenantSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "65000"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "65100"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "description", "Test ASN range with full options"),

					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "rir"),

					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "tenant"),
				),
			},
		},
	})

}

func TestAccASNRangeResource_update(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-asn-range-upd")

	slug := testutil.RandomSlug("tf-test-asn-range-upd")

	rirName := testutil.RandomName("tf-test-rir")

	rirSlug := testutil.RandomSlug("tf-test-rir")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterASNRangeCleanup(name)

	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckASNRangeDestroy,

			testutil.CheckRIRDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "64512"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "64612"),
				),
			},

			{

				Config: testAccASNRangeResourceConfig_updated(name, slug, rirName, rirSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "64512"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "64700"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "description", "Updated description"),
				),
			},
		},
	})

}

func testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug string) string {

	return fmt.Sprintf(`



resource "netbox_rir" "test" {



  name = %q



  slug = %q



}



resource "netbox_asn_range" "test" {



  name  = %q



  slug  = %q



  rir   = netbox_rir.test.id



  start = "64512"



  end   = "64612"



}



`, rirName, rirSlug, name, slug)

}

func testAccASNRangeResourceConfig_full(name, slug, rirName, rirSlug, tenantName, tenantSlug string) string {

	return fmt.Sprintf(`



resource "netbox_rir" "test" {



  name = %q



  slug = %q



}







resource "netbox_tenant" "test" {



  name = %q



  slug = %q



}







resource "netbox_asn_range" "test" {



  name        = %q



  slug        = %q



  rir         = netbox_rir.test.id



  start       = "65000"



  end         = "65100"



  tenant      = netbox_tenant.test.id



  description = "Test ASN range with full options"



}



`, rirName, rirSlug, tenantName, tenantSlug, name, slug)

}

func testAccASNRangeResourceConfig_updated(name, slug, rirName, rirSlug string) string {

	return fmt.Sprintf(`



resource "netbox_rir" "test" {



  name = %q



  slug = %q



}



resource "netbox_asn_range" "test" {



  name        = %q



  slug        = %q



  rir         = netbox_rir.test.id



  start       = "64512"



  end         = "64700"



  description = "Updated description"



}



`, rirName, rirSlug, name, slug)

}

func TestAccASNRangeResource_import(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	name := testutil.RandomName("tf-test-asn-range")

	slug := testutil.RandomSlug("tf-test-asn-range")

	rirName := testutil.RandomName("tf-test-rir")

	rirSlug := testutil.RandomSlug("tf-test-rir")

	// Register cleanup to ensure resources are deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterASNRangeCleanup(name)

	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckASNRangeDestroy,

			testutil.CheckRIRDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "64512"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "64612"),

					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "rir"),
				),
			},

			{

				ResourceName: "netbox_asn_range.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"rir"},
			},
		},
	})

}

func TestAccConsistency_ASNRange(t *testing.T) {
	rangeName := testutil.RandomName("asn-range")

	rangeSlug := testutil.RandomSlug("asn-range")

	rirName := testutil.RandomName("rir")

	rirSlug := testutil.RandomSlug("rir")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccASNRangeConsistencyConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", rangeName),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "rir", rirSlug),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "tenant", tenantName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccASNRangeConsistencyConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug),
			},
		},
	})

}

func testAccASNRangeConsistencyConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug string) string {

	return fmt.Sprintf(`



resource "netbox_rir" "test" {



  name = "%[3]s"



  slug = "%[4]s"



}







resource "netbox_tenant" "test" {



  name = "%[5]s"



  slug = "%[6]s"



}







resource "netbox_asn_range" "test" {



  name = "%[1]s"



  slug = "%[2]s"



  rir = netbox_rir.test.slug



  tenant = netbox_tenant.test.name



  start = 65000



  end = 65100



}



`, rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug)

}

// TestAccConsistency_ASNRange_LiteralNames tests that reference attributes specified as literal string names

// are preserved and do not cause drift when the API returns numeric IDs.

func TestAccConsistency_ASNRange_LiteralNames(t *testing.T) {
	rangeName := testutil.RandomName("asn-range")

	rangeSlug := testutil.RandomSlug("asn-range")

	rirName := testutil.RandomName("rir")

	rirSlug := testutil.RandomSlug("rir")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccASNRangeConsistencyLiteralNamesConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", rangeName),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "rir", rirSlug),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "tenant", tenantName),
				),
			},

			{

				// Critical: Verify no drift when refreshing state

				PlanOnly: true,

				Config: testAccASNRangeConsistencyLiteralNamesConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug),
			},
		},
	})

}

func testAccASNRangeConsistencyLiteralNamesConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug string) string {

	return fmt.Sprintf(`



resource "netbox_rir" "test" {



  name = "%[3]s"



  slug = "%[4]s"



}







resource "netbox_tenant" "test" {



  name = "%[5]s"



  slug = "%[6]s"



}







resource "netbox_asn_range" "test" {



  name = "%[1]s"



  slug = "%[2]s"







  # Use literal string names to mimic existing user state



  rir = "%[4]s"



  tenant = "%[5]s"



  start = 65000



  end = 65100







  depends_on = [netbox_rir.test, netbox_tenant.test]



}



`, rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug)

}
