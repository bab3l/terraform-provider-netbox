package resources_acceptance_tests

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSiteASNAssignmentResource_basic(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")
	asn := int64(acctest.RandIntRange(1000000000, 2000000000))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterASNCleanup(asn)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteASNAssignmentResourceConfig_basic(siteName, siteSlug, rirName, rirSlug, asn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site_asn_assignment.test", "id"),
					resource.TestCheckResourceAttrPair("netbox_site_asn_assignment.test", "site", "netbox_site.test", "slug"),
					resource.TestCheckResourceAttrPair("netbox_site_asn_assignment.test", "asn", "netbox_asn.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_site_asn_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"site",
				},
			},
		},
	})
}

func TestAccSiteASNAssignmentResource_full(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-full")
	siteSlug := testutil.RandomSlug("tf-test-site-full")
	rirName := testutil.RandomName("tf-test-rir-full")
	rirSlug := testutil.RandomSlug("tf-test-rir-full")
	asn := int64(acctest.RandIntRange(1000000000, 2000000000))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterASNCleanup(asn)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteASNAssignmentResourceConfig_basic(siteName, siteSlug, rirName, rirSlug, asn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site_asn_assignment.test", "id"),
					resource.TestCheckResourceAttrPair("netbox_site_asn_assignment.test", "site", "netbox_site.test", "slug"),
					resource.TestCheckResourceAttrPair("netbox_site_asn_assignment.test", "asn", "netbox_asn.test", "id"),
				),
			},
		},
	})
}

func TestAccSiteASNAssignmentResource_update(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-update")
	siteSlug := testutil.RandomSlug("tf-test-site-update")
	rirName := testutil.RandomName("tf-test-rir-update")
	rirSlug := testutil.RandomSlug("tf-test-rir-update")
	asn1 := int64(acctest.RandIntRange(1000000000, 1499999999))
	asn2 := int64(acctest.RandIntRange(1500000000, 2000000000))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterASNCleanup(asn1)
	cleanup.RegisterASNCleanup(asn2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteASNAssignmentResourceConfig_twoASNs(siteName, siteSlug, rirName, rirSlug, asn1, asn2, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("netbox_site_asn_assignment.test", "asn", "netbox_asn.first", "id"),
				),
			},
			{
				Config: testAccSiteASNAssignmentResourceConfig_twoASNs(siteName, siteSlug, rirName, rirSlug, asn1, asn2, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("netbox_site_asn_assignment.test", "asn", "netbox_asn.second", "id"),
				),
			},
		},
	})
}

func TestAccSiteASNAssignmentResource_import(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-import")
	siteSlug := testutil.RandomSlug("tf-test-site-import")
	rirName := testutil.RandomName("tf-test-rir-import")
	rirSlug := testutil.RandomSlug("tf-test-rir-import")
	asn := int64(acctest.RandIntRange(1000000000, 2000000000))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterASNCleanup(asn)

	testutil.RunImportTest(t, testutil.ImportTestConfig{
		ResourceName: "netbox_site_asn_assignment",
		Config: func() string {
			return testAccSiteASNAssignmentResourceConfig_basic(siteName, siteSlug, rirName, rirSlug, asn)
		},
		ImportStateVerifyIgnore: []string{
			"site",
		},
		CheckDestroy: testutil.CheckSiteDestroy,
	})
}

func TestAccSiteASNAssignmentResource_externalDeletion(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-ext")
	siteSlug := testutil.RandomSlug("tf-test-site-ext")
	rirName := testutil.RandomName("tf-test-rir-ext")
	rirSlug := testutil.RandomSlug("tf-test-rir-ext")
	asn := int64(acctest.RandIntRange(1000000000, 2000000000))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterASNCleanup(asn)

	testutil.RunExternalDeletionTest(t, testutil.ExternalDeletionTestConfig{
		ResourceName: "netbox_site_asn_assignment",
		Config: func() string {
			return testAccSiteASNAssignmentResourceConfig_basic(siteName, siteSlug, rirName, rirSlug, asn)
		},
		DeleteFunc: func(ctx context.Context, id string) error {
			parts := strings.Split(id, ":")
			if len(parts) != 2 {
				return fmt.Errorf("invalid ID format: %s", id)
			}
			siteID, err := utils.ParseID(parts[0])
			if err != nil {
				return err
			}
			asnID, err := utils.ParseID(parts[1])
			if err != nil {
				return err
			}

			client, err := testutil.GetSharedClient()
			if err != nil {
				return err
			}
			site, httpResp, err := client.DcimAPI.DcimSitesRetrieve(ctx, siteID).Execute()
			defer utils.CloseResponseBody(httpResp)
			if err != nil {
				return err
			}

			asnIDs := make([]int32, 0, len(site.GetAsns()))
			for _, asnObj := range site.GetAsns() {
				if asnObj.GetId() == asnID {
					continue
				}
				asnIDs = append(asnIDs, asnObj.GetId())
			}

			siteRequest := netbox.WritableSiteRequest{
				Name: site.GetName(),
				Slug: site.GetSlug(),
			}
			siteRequest.SetAsns(asnIDs)

			_, updateResp, updateErr := client.DcimAPI.DcimSitesUpdate(ctx, siteID).WritableSiteRequest(siteRequest).Execute()
			defer utils.CloseResponseBody(updateResp)
			return updateErr
		},
		CheckDestroy: testutil.CheckSiteDestroy,
	})
}

func TestAccSiteASNAssignmentResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-remove")
	siteSlug := testutil.RandomSlug("tf-test-site-remove")
	rirName := testutil.RandomName("tf-test-rir-remove")
	rirSlug := testutil.RandomSlug("tf-test-rir-remove")
	asn := int64(acctest.RandIntRange(1000000000, 2000000000))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterASNCleanup(asn)

	baseConfig := func() string {
		return testAccSiteASNAssignmentResourceConfig_basic(siteName, siteSlug, rirName, rirSlug, asn)
	}

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_site_asn_assignment",
		BaseConfig:   baseConfig,
		ConfigWithFields: func() string {
			return baseConfig()
		},
		OptionalFields: map[string]string{},
		RequiredFields: map[string]string{
			"site": siteSlug,
		},
		CheckDestroy: testutil.CheckSiteDestroy,
	})
}

func testAccSiteASNAssignmentResourceConfig_basic(siteName, siteSlug, rirName, rirSlug string, asn int64) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn" "test" {
  asn = %d
  rir = netbox_rir.test.id
}

resource "netbox_site_asn_assignment" "test" {
  site = netbox_site.test.slug
  asn  = netbox_asn.test.id
}
`, siteName, siteSlug, rirName, rirSlug, asn)
}

func testAccSiteASNAssignmentResourceConfig_twoASNs(siteName, siteSlug, rirName, rirSlug string, asn1, asn2 int64, useSecond bool) string {
	assignmentASN := "netbox_asn.first.id"
	if useSecond {
		assignmentASN = "netbox_asn.second.id"
	}

	return fmt.Sprintf(`
resource "netbox_site" "test" {
	name   = %q
	slug   = %q
	status = "active"
}

resource "netbox_rir" "test" {
	name = %q
	slug = %q
}

resource "netbox_asn" "first" {
	asn = %d
	rir = netbox_rir.test.id
}

resource "netbox_asn" "second" {
	asn = %d
	rir = netbox_rir.test.id
}

resource "netbox_site_asn_assignment" "test" {
	site = netbox_site.test.slug
	asn  = %s
}
`, siteName, siteSlug, rirName, rirSlug, asn1, asn2, assignmentASN)
}
