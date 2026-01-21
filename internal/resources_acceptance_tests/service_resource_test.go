package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceResource_basic(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	mfgName := testutil.RandomName("tf-test-mfg")
	mfgSlug := testutil.RandomSlug("tf-test-mfg")
	dtModel := testutil.RandomName("tf-test-dt")
	dtSlug := testutil.RandomSlug("tf-test-dt")
	roleName := testutil.RandomName("tf-test-role")
	roleSlug := testutil.RandomSlug("tf-test-role")
	deviceName := testutil.RandomName("tf-test-device")
	serviceName := testutil.RandomName("tf-test-svc")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_service.test", "id"),
					resource.TestCheckResourceAttr("netbox_service.test", "name", serviceName),
					resource.TestCheckResourceAttr("netbox_service.test", "protocol", "tcp"),
				),
			},
			{
				Config:   testAccServiceResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName),
				PlanOnly: true,
			},
			{
				ResourceName:            "netbox_service.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device"},
				Check: resource.ComposeTestCheckFunc(
					testutil.ReferenceFieldCheck("netbox_service.test", "device"),
					testutil.ReferenceFieldCheck("netbox_service.test", "virtual_machine"),
				),
			},
			{
				Config:   testAccServiceResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccServiceResource_full(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-full")
	siteSlug := testutil.RandomSlug("tf-test-site-full")
	mfgName := testutil.RandomName("tf-test-mfg-full")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-full")
	dtModel := testutil.RandomName("tf-test-dt-full")
	dtSlug := testutil.RandomSlug("tf-test-dt-full")
	roleName := testutil.RandomName("tf-test-role-full")
	roleSlug := testutil.RandomSlug("tf-test-role-full")
	deviceName := testutil.RandomName("tf-test-device-full")
	serviceName := testutil.RandomName("tf-test-svc-full")
	description := testutil.RandomName("description")
	updatedDescription := "Updated service description"
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, description, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_service.test", "id"),
					resource.TestCheckResourceAttr("netbox_service.test", "name", serviceName),
					resource.TestCheckResourceAttr("netbox_service.test", "protocol", "tcp"),
					resource.TestCheckResourceAttr("netbox_service.test", "description", description),
					resource.TestCheckResourceAttr("netbox_service.test", "tags.#", "2"),
				),
			},
			{
				Config:   testAccServiceResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, description, tagName1, tagSlug1, tagName2, tagSlug2),
				PlanOnly: true,
			},
			{
				Config: testAccServiceResourceConfig_fullUpdate(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, updatedDescription, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_service.test", "comments", "Updated comments for service"),
				),
			},
			{
				Config:   testAccServiceResourceConfig_fullUpdate(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, updatedDescription, tagName1, tagSlug1, tagName2, tagSlug2),
				PlanOnly: true,
			},
		},
	})
}

func TestAccServiceResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-tags")
	siteSlug := testutil.RandomSlug("tf-test-site-tags")
	mfgName := testutil.RandomName("tf-test-mfg-tags")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-tags")
	dtModel := testutil.RandomName("tf-test-dt-tags")
	dtSlug := testutil.RandomSlug("tf-test-dt-tags")
	roleName := testutil.RandomName("tf-test-role-tags")
	roleSlug := testutil.RandomSlug("tf-test-role-tags")
	deviceName := testutil.RandomName("tf-test-device-tags")
	serviceName := testutil.RandomName("tf-test-svc-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceResourceConfig_tags(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_service.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_service.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccServiceResourceConfig_tags(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_service.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_service.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccServiceResourceConfig_tags(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_service.test", "tags.*", tag3Slug),
				),
			},
			{
				Config: testAccServiceResourceConfig_tags(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccServiceResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-tag-order")
	siteSlug := testutil.RandomSlug("tf-test-site-tag-order")
	mfgName := testutil.RandomName("tf-test-mfg-tag-order")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-tag-order")
	dtModel := testutil.RandomName("tf-test-dt-tag-order")
	dtSlug := testutil.RandomSlug("tf-test-dt-tag-order")
	roleName := testutil.RandomName("tf-test-role-tag-order")
	roleSlug := testutil.RandomSlug("tf-test-role-tag-order")
	deviceName := testutil.RandomName("tf-test-device-tag-order")
	serviceName := testutil.RandomName("tf-test-svc-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceResourceConfig_tagsOrder(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_service.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_service.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccServiceResourceConfig_tagsOrder(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_service.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_service.test", "tags.*", tag2Slug),
				),
			},
		},
	})
}

func TestAccConsistency_Service(t *testing.T) {
	t.Parallel()

	serviceName := testutil.RandomName("service")
	vmName := testutil.RandomName("vm")
	clusterName := testutil.RandomName("cluster")
	clusterTypeName := testutil.RandomName("cluster-type")
	clusterTypeSlug := testutil.RandomSlug("cluster-type")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceConsistencyConfig(serviceName, vmName, clusterName, clusterTypeName, clusterTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service.test", "name", serviceName),
					resource.TestCheckResourceAttr("netbox_service.test", "virtual_machine", vmName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccServiceConsistencyConfig(serviceName, vmName, clusterName, clusterTypeName, clusterTypeSlug),
			},
		},
	})
}

func TestAccConsistency_Service_LiteralNames(t *testing.T) {
	t.Parallel()

	serviceName := testutil.RandomName("service")
	vmName := testutil.RandomName("vm")
	clusterName := testutil.RandomName("cluster")
	clusterTypeName := testutil.RandomName("cluster-type")
	clusterTypeSlug := testutil.RandomSlug("cluster-type")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceConsistencyLiteralNamesConfig(serviceName, vmName, clusterName, clusterTypeName, clusterTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service.test", "name", serviceName),
					resource.TestCheckResourceAttr("netbox_service.test", "virtual_machine", vmName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccServiceConsistencyLiteralNamesConfig(serviceName, vmName, clusterName, clusterTypeName, clusterTypeSlug),
			},
		},
	})
}

func TestAccServiceResource_update(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-upd")
	siteSlug := testutil.RandomSlug("tf-test-site-upd")
	mfgName := testutil.RandomName("tf-test-mfg-upd")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-upd")
	dtModel := testutil.RandomName("tf-test-dt-upd")
	dtSlug := testutil.RandomSlug("tf-test-dt-upd")
	roleName := testutil.RandomName("tf-test-role-upd")
	roleSlug := testutil.RandomSlug("tf-test-role-upd")
	deviceName := testutil.RandomName("tf-test-device-upd")
	serviceName := testutil.RandomName("tf-test-svc-upd")
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, testutil.Description1, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_service.test", "id"),
					resource.TestCheckResourceAttr("netbox_service.test", "name", serviceName),
					resource.TestCheckResourceAttr("netbox_service.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccServiceResourceConfig_fullUpdate(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, testutil.Description2, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccServiceResource_externalDeletion(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-ext")
	siteSlug := testutil.RandomSlug("tf-test-site-ext")
	mfgName := testutil.RandomName("tf-test-mfg-ext")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-ext")
	dtModel := testutil.RandomName("tf-test-dt-ext")
	dtSlug := testutil.RandomSlug("tf-test-dt-ext")
	roleName := testutil.RandomName("tf-test-role-ext")
	roleSlug := testutil.RandomSlug("tf-test-role-ext")
	deviceName := testutil.RandomName("tf-test-device-ext")
	serviceName := testutil.RandomName("tf-test-svc-ext")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_service.test", "id"),
					resource.TestCheckResourceAttr("netbox_service.test", "name", serviceName),
					resource.TestCheckResourceAttr("netbox_service.test", "protocol", "tcp"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.IpamAPI.IpamServicesList(context.Background()).Name([]string{serviceName}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find service for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamServicesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete service: %v", err)
					}
					t.Logf("Successfully externally deleted service with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccServiceResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_device_role" "test" {
  name  = %q
  slug  = %q
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_service" "test" {
  device   = netbox_device.test.id
  name     = %q
  protocol = "tcp"
  ports    = [22]
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName)
}

func testAccServiceResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, description, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = %[10]q
	slug = %[11]q
}

resource "netbox_tag" "tag2" {
	name = %[12]q
	slug = %[13]q
}

resource "netbox_site" "test" {
	name   = %[1]q
	slug   = %[2]q
  status = "active"
}

resource "netbox_manufacturer" "test" {
	name = %[3]q
	slug = %[4]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
	model        = %[5]q
	slug         = %[6]q
}

resource "netbox_device_role" "test" {
	name  = %[7]q
	slug  = %[8]q
  color = "aa1409"
}

resource "netbox_device" "test" {
	name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_service" "test" {
  device      = netbox_device.test.id
	name        = %[14]q
  protocol    = "tcp"
  ports       = [22, 443]
	description = %[15]q

	tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, tagName1, tagSlug1, tagName2, tagSlug2, serviceName, description)
}

func testAccServiceResourceConfig_fullUpdate(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, description, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = %[10]q
	slug = %[11]q
}

resource "netbox_tag" "tag2" {
	name = %[12]q
	slug = %[13]q
}

resource "netbox_site" "test" {
	name   = %[1]q
	slug   = %[2]q
	status = "active"
}

resource "netbox_manufacturer" "test" {
	name = %[3]q
	slug = %[4]q
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = %[5]q
	slug         = %[6]q
}

resource "netbox_device_role" "test" {
	name  = %[7]q
	slug  = %[8]q
	color = "aa1409"
}

resource "netbox_device" "test" {
	name        = %[9]q
	device_type = netbox_device_type.test.id
	role        = netbox_device_role.test.id
	site        = netbox_site.test.id
}

resource "netbox_service" "test" {
	device      = netbox_device.test.id
	name        = %[14]q
	protocol    = "tcp"
	ports       = [22, 443]
	description = %[15]q
	comments    = "Updated comments for service"

	tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, tagName1, tagSlug1, tagName2, tagSlug2, serviceName, description)
}

func testAccServiceResourceConfig_tags(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
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
	name = "Tag1-%[11]s"
	slug = %[11]q
}

resource "netbox_tag" "tag2" {
	name = "Tag2-%[12]s"
	slug = %[12]q
}

resource "netbox_tag" "tag3" {
	name = "Tag3-%[13]s"
	slug = %[13]q
}

resource "netbox_site" "test" {
	name   = %[1]q
	slug   = %[2]q
	status = "active"
}

resource "netbox_manufacturer" "test" {
	name = %[3]q
	slug = %[4]q
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = %[5]q
	slug         = %[6]q
}

resource "netbox_device_role" "test" {
	name  = %[7]q
	slug  = %[8]q
	color = "aa1409"
}

resource "netbox_device" "test" {
	name        = %[9]q
	device_type = netbox_device_type.test.id
	role        = netbox_device_role.test.id
	site        = netbox_site.test.id
}

resource "netbox_service" "test" {
	device   = netbox_device.test.id
	name     = %[10]q
	protocol = "tcp"
	ports    = [22]
	%[14]s
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccServiceResourceConfig_tagsOrder(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleSlugReversed
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = "Tag1-%[11]s"
	slug = %[11]q
}

resource "netbox_tag" "tag2" {
	name = "Tag2-%[12]s"
	slug = %[12]q
}

resource "netbox_site" "test" {
	name   = %[1]q
	slug   = %[2]q
	status = "active"
}

resource "netbox_manufacturer" "test" {
	name = %[3]q
	slug = %[4]q
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = %[5]q
	slug         = %[6]q
}

resource "netbox_device_role" "test" {
	name  = %[7]q
	slug  = %[8]q
	color = "aa1409"
}

resource "netbox_device" "test" {
	name        = %[9]q
	device_type = netbox_device_type.test.id
	role        = netbox_device_role.test.id
	site        = netbox_site.test.id
}

resource "netbox_service" "test" {
	device   = netbox_device.test.id
	name     = %[10]q
	protocol = "tcp"
	ports    = [22]
	%[13]s
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, tag1Slug, tag2Slug, tagsConfig)
}

func testAccServiceConsistencyConfig(serviceName, vmName, clusterName, clusterTypeName, clusterTypeSlug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_cluster" "test" {
  name = "%[3]s"
  type = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
  name = "%[2]s"
  cluster = netbox_cluster.test.id
}

resource "netbox_service" "test" {
  name = "%[1]s"
  virtual_machine = netbox_virtual_machine.test.name
  ports = [80]
  protocol = "tcp"
}
`, serviceName, vmName, clusterName, clusterTypeName, clusterTypeSlug)
}

func testAccServiceConsistencyLiteralNamesConfig(serviceName, vmName, clusterName, clusterTypeName, clusterTypeSlug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_cluster" "test" {
  name = "%[3]s"
  type = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
  name = "%[2]s"
  cluster = netbox_cluster.test.id
}

resource "netbox_service" "test" {
  name = "%[1]s"
  # Use literal string name to mimic existing user state
  virtual_machine = "%[2]s"
  ports = [80]
  protocol = "tcp"

  depends_on = [netbox_virtual_machine.test]
}
`, serviceName, vmName, clusterName, clusterTypeName, clusterTypeSlug)
}

func TestAccServiceResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	// Unique names and slugs
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	mfgName := testutil.RandomName("tf-test-mfg")
	mfgSlug := testutil.RandomSlug("tf-test-mfg")
	dtModel := testutil.RandomName("tf-test-dt")
	dtSlug := testutil.RandomSlug("tf-test-dt")
	roleName := testutil.RandomName("tf-test-role")
	roleSlug := testutil.RandomSlug("tf-test-role")
	deviceName := testutil.RandomName("tf-test-device")
	serviceName := testutil.RandomName("tf-test-svc")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	// Common dependency config
	depsConfig := fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_device_role" "test" {
  name  = %q
  slug  = %q
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName)

	config := testutil.MultiFieldOptionalTestConfig{
		ResourceType: "netbox_service",
		ResourceName: "netbox_service",
		ConfigWithFields: func() string {
			return depsConfig + fmt.Sprintf(`
resource "netbox_service" "test" {
  device      = netbox_device.test.id
  name        = %q
  protocol    = "tcp"
  ports       = [80]
  description = "Test Description"
  comments    = "Test Comments"
}
`, serviceName)
		},
		BaseConfig: func() string {
			return depsConfig + fmt.Sprintf(`
resource "netbox_service" "test" {
  device      = netbox_device.test.id
  name        = %q
  protocol    = "tcp"
  ports       = [80]
}
`, serviceName)
		},
		OptionalFields: map[string]string{
			"description": "Test Description",
			"comments":    "Test Comments",
		},
	}

	testutil.TestRemoveOptionalFields(t, config)
}
func TestAccServiceResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_service",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_site" "test" {
  name = "test-site"
  slug = "test-site"
}

resource "netbox_device_role" "test" {
  name = "test-role"
  slug = "test-role"
}

resource "netbox_device_type" "test" {
  model = "test-type"
  slug  = "test-type"
}

resource "netbox_device" "test" {
  name        = "test-device"
  site        = netbox_site.test.id
  device_role = netbox_device_role.test.id
  device_type = netbox_device_type.test.id
}

resource "netbox_service" "test" {
  # name missing
  device   = netbox_device.test.id
  protocol = "tcp"
  ports    = [80]
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_protocol": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_site" "test" {
  name = "test-site"
  slug = "test-site"
}

resource "netbox_device_role" "test" {
  name = "test-role"
  slug = "test-role"
}

resource "netbox_device_type" "test" {
  model = "test-type"
  slug  = "test-type"
}

resource "netbox_device" "test" {
  name        = "test-device"
  site        = netbox_site.test.id
  device_role = netbox_device_role.test.id
  device_type = netbox_device_type.test.id
}

resource "netbox_service" "test" {
  name   = "test-service"
  device = netbox_device.test.id
  # protocol missing
  ports  = [80]
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_ports": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_site" "test" {
  name = "test-site"
  slug = "test-site"
}

resource "netbox_device_role" "test" {
  name = "test-role"
  slug = "test-role"
}

resource "netbox_device_type" "test" {
  model = "test-type"
  slug  = "test-type"
}

resource "netbox_device" "test" {
  name        = "test-device"
  site        = netbox_site.test.id
  device_role = netbox_device_role.test.id
  device_type = netbox_device_type.test.id
}

resource "netbox_service" "test" {
  name     = "test-service"
  device   = netbox_device.test.id
  protocol = "tcp"
  # ports missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
