package resources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestInventoryItemTemplateResource(t *testing.T) {

	t.Parallel()

	r := resources.NewInventoryItemTemplateResource()

	if r == nil {

		t.Fatal("Expected non-nil InventoryItemTemplate resource")
	}
}

func TestInventoryItemTemplateResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewInventoryItemTemplateResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")
	}

	requiredAttrs := []string{"device_type", "name"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected required attribute %s to exist in schema", attr)
		}
	}

	computedAttrs := []string{"id"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}

	optionalAttrs := []string{"parent", "label", "role", "manufacturer", "part_id", "description", "component_type", "component_id"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestInventoryItemTemplateResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewInventoryItemTemplateResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_inventory_item_template"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestInventoryItemTemplateResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewInventoryItemTemplateResource().(*resources.InventoryItemTemplateResource)

	configureRequest := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)
	}

	configureRequest.ProviderData = invalidProviderData

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with incorrect provider data")
	}
}

func TestAccInventoryItemTemplateResource_basic(t *testing.T) {

	mfgName := testutil.RandomName("tf-test-mfg")

	mfgSlug := testutil.RandomSlug("tf-test-mfg")

	dtModel := testutil.RandomName("tf-test-dt")

	dtSlug := testutil.RandomSlug("tf-test-dt")

	templateName := testutil.RandomName("tf-test-iit")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterManufacturerCleanup(mfgSlug)

	cleanup.RegisterDeviceTypeCleanup(dtSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckInventoryItemTemplateDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccInventoryItemTemplateResourceConfig_basic(mfgName, mfgSlug, dtModel, dtSlug, templateName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_inventory_item_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_inventory_item_template.test", "name", templateName),
				),
			},

			// ImportState test

			{

				ResourceName: "netbox_inventory_item_template.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})
}

func TestAccInventoryItemTemplateResource_full(t *testing.T) {

	mfgName := testutil.RandomName("tf-test-mfg")

	mfgSlug := testutil.RandomSlug("tf-test-mfg")

	dtModel := testutil.RandomName("tf-test-dt")

	dtSlug := testutil.RandomSlug("tf-test-dt")

	templateName := testutil.RandomName("tf-test-iit")

	label := "Component Label"

	partID := "PART-12345"

	description := "Test inventory item template"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterManufacturerCleanup(mfgSlug)

	cleanup.RegisterDeviceTypeCleanup(dtSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckInventoryItemTemplateDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccInventoryItemTemplateResourceConfig_full(mfgName, mfgSlug, dtModel, dtSlug, templateName, label, partID, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_inventory_item_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_inventory_item_template.test", "name", templateName),

					resource.TestCheckResourceAttr("netbox_inventory_item_template.test", "label", label),

					resource.TestCheckResourceAttr("netbox_inventory_item_template.test", "part_id", partID),

					resource.TestCheckResourceAttr("netbox_inventory_item_template.test", "description", description),
				),
			},
		},
	})
}

func TestAccInventoryItemTemplateResource_update(t *testing.T) {

	mfgName := testutil.RandomName("tf-test-mfg")

	mfgSlug := testutil.RandomSlug("tf-test-mfg")

	dtModel := testutil.RandomName("tf-test-dt")

	dtSlug := testutil.RandomSlug("tf-test-dt")

	templateName := testutil.RandomName("tf-test-iit")

	const description1 = "Initial description"

	const description2 = "Updated description"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterManufacturerCleanup(mfgSlug)

	cleanup.RegisterDeviceTypeCleanup(dtSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckInventoryItemTemplateDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccInventoryItemTemplateResourceConfig_full(mfgName, mfgSlug, dtModel, dtSlug, templateName, "", "", description1),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_inventory_item_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_inventory_item_template.test", "description", description1),
				),
			},

			{

				Config: testAccInventoryItemTemplateResourceConfig_full(mfgName, mfgSlug, dtModel, dtSlug, templateName, "", "", description2),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_inventory_item_template.test", "description", description2),
				),
			},
		},
	})
}

func TestAccInventoryItemTemplateResource_withParent(t *testing.T) {

	mfgName := testutil.RandomName("tf-test-mfg")

	mfgSlug := testutil.RandomSlug("tf-test-mfg")

	dtModel := testutil.RandomName("tf-test-dt")

	dtSlug := testutil.RandomSlug("tf-test-dt")

	parentName := testutil.RandomName("tf-test-parent")

	childName := testutil.RandomName("tf-test-child")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterManufacturerCleanup(mfgSlug)

	cleanup.RegisterDeviceTypeCleanup(dtSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckInventoryItemTemplateDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccInventoryItemTemplateResourceConfig_withParent(mfgName, mfgSlug, dtModel, dtSlug, parentName, childName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_inventory_item_template.parent", "id"),

					resource.TestCheckResourceAttr("netbox_inventory_item_template.parent", "name", parentName),

					resource.TestCheckResourceAttrSet("netbox_inventory_item_template.child", "id"),

					resource.TestCheckResourceAttr("netbox_inventory_item_template.child", "name", childName),

					resource.TestCheckResourceAttrPair("netbox_inventory_item_template.child", "parent", "netbox_inventory_item_template.parent", "id"),
				),
			},
		},
	})
}

func testAccInventoryItemTemplateResourceConfig_basic(mfgName, mfgSlug, dtModel, dtSlug, templateName string) string {

	return fmt.Sprintf(`

provider "netbox" {}

resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  model        = %[3]q
  slug         = %[4]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_inventory_item_template" "test" {
  name        = %[5]q

  device_type = netbox_device_type.test.id
}

`, mfgName, mfgSlug, dtModel, dtSlug, templateName)
}

func testAccInventoryItemTemplateResourceConfig_full(mfgName, mfgSlug, dtModel, dtSlug, templateName, label, partID, description string) string {

	labelAttr := ""

	if label != "" {

		labelAttr = fmt.Sprintf(`label       = %q`, label)
	}

	partIDAttr := ""

	if partID != "" {

		partIDAttr = fmt.Sprintf(`part_id     = %q`, partID)
	}

	descAttr := ""

	if description != "" {

		descAttr = fmt.Sprintf(`description = %q`, description)
	}

	return fmt.Sprintf(`

provider "netbox" {}

resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  model        = %[3]q
  slug         = %[4]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_inventory_item_template" "test" {
  name        = %[5]q

  device_type = netbox_device_type.test.id

  %[6]s

  %[7]s

  %[8]s
}

`, mfgName, mfgSlug, dtModel, dtSlug, templateName, labelAttr, partIDAttr, descAttr)
}

func testAccInventoryItemTemplateResourceConfig_withParent(mfgName, mfgSlug, dtModel, dtSlug, parentName, childName string) string {

	return fmt.Sprintf(`

provider "netbox" {}

resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  model        = %[3]q
  slug         = %[4]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_inventory_item_template" "parent" {
  name        = %[5]q

  device_type = netbox_device_type.test.id
}

resource "netbox_inventory_item_template" "child" {
  name        = %[6]q

  device_type = netbox_device_type.test.id
  parent      = netbox_inventory_item_template.parent.id
}

`, mfgName, mfgSlug, dtModel, dtSlug, parentName, childName)
}
