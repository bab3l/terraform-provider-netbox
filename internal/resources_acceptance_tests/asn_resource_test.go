package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccASNResource_basic(t *testing.T) {

	rirName := testutil.RandomName("tf-test-rir")

	rirSlug := testutil.RandomSlug("tf-test-rir")

	// Generate a random ASN in the private range (64512-65534)

	asn := int64(acctest.RandIntRange(64512, 65534))

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccASNResourceConfig_basic(rirName, rirSlug, asn),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_asn.test", "id"),

					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
				),
			},

			{

				ResourceName: "netbox_asn.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"rir"},
			},
		},
	})

}

func TestAccASNResource_full(t *testing.T) {

	rirName := testutil.RandomName("tf-test-rir")

	rirSlug := testutil.RandomSlug("tf-test-rir")

	// Generate a random ASN in the private range (64512-65534)

	asn := int64(acctest.RandIntRange(64512, 65534))

	description := "Test ASN with all fields"

	updatedDescription := "Updated ASN description"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccASNResourceConfig_full(rirName, rirSlug, asn, description, testutil.Comments),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_asn.test", "id"),

					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),

					resource.TestCheckResourceAttr("netbox_asn.test", "description", description),

					resource.TestCheckResourceAttr("netbox_asn.test", "comments", testutil.Comments),

					resource.TestCheckResourceAttrSet("netbox_asn.test", "rir"),
				),
			},

			{

				Config: testAccASNResourceConfig_full(rirName, rirSlug, asn, updatedDescription, testutil.Comments),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_asn.test", "description", updatedDescription),
				),
			},
		},
	})

}

func testAccASNResourceConfig_basic(rirName, rirSlug string, asn int64) string {

	return fmt.Sprintf(`



resource "netbox_rir" "test" {



  name = %q



  slug = %q



}







resource "netbox_asn" "test" {



  asn = %d



  rir = netbox_rir.test.id



}



`, rirName, rirSlug, asn)

}

func testAccASNResourceConfig_full(rirName, rirSlug string, asn int64, description, comments string) string {

	return fmt.Sprintf(`



resource "netbox_rir" "test" {



  name = %q



  slug = %q



}







resource "netbox_asn" "test" {



  asn         = %d



  rir         = netbox_rir.test.id



  description = %q



  comments    = %q



}



`, rirName, rirSlug, asn, description, comments)

}

// TestAccConsistency_ASN_LiteralNames tests that reference attributes specified as literal string names

// are preserved and do not cause drift when the API returns numeric IDs.

func TestAccConsistency_ASN_LiteralNames(t *testing.T) {

	t.Parallel()

	asn := int64(65100)

	rirName := testutil.RandomName("rir")

	rirSlug := testutil.RandomSlug("rir")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccASNConsistencyLiteralNamesConfig(asn, rirName, rirSlug, tenantName, tenantSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),

					resource.TestCheckResourceAttr("netbox_asn.test", "rir", rirSlug),

					resource.TestCheckResourceAttr("netbox_asn.test", "tenant", tenantName),
				),
			},

			{

				// Critical: Verify no drift when refreshing state

				PlanOnly: true,

				Config: testAccASNConsistencyLiteralNamesConfig(asn, rirName, rirSlug, tenantName, tenantSlug),
			},
		},
	})

}

func testAccASNConsistencyLiteralNamesConfig(asn int64, rirName, rirSlug, tenantName, tenantSlug string) string {

	return fmt.Sprintf(`



resource "netbox_rir" "test" {



  name = "%[2]s"



  slug = "%[3]s"



}







resource "netbox_tenant" "test" {



  name = "%[4]s"



  slug = "%[5]s"



}







resource "netbox_asn" "test" {



  asn = %[1]d



  # Use literal string names to mimic existing user state



  rir = "%[3]s"



  tenant = "%[4]s"







  depends_on = [netbox_rir.test, netbox_tenant.test]



}



`, asn, rirName, rirSlug, tenantName, tenantSlug)

}
