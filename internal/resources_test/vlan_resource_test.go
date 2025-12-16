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

func TestVLANResource(t *testing.T) {

	t.Parallel()

	r := resources.NewVLANResource()

	if r == nil {

		t.Fatal("Expected non-nil VLAN resource")

	}

}

func TestVLANResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewVLANResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"name", "vid"}

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

	optionalAttrs := []string{"status", "site", "group", "tenant", "role", "description", "comments"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestVLANResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewVLANResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_vlan"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestVLANResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewVLANResource().(*resources.VLANResource)

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

func TestAccVLANResource_basic(t *testing.T) {

	name := testutil.RandomName("tf-test-vlan")

	vid := testutil.RandomVID()

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterVLANCleanup(vid)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVLANDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccVLANResourceConfig_basic(name, vid),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),

					resource.TestCheckResourceAttr("netbox_vlan.test", "name", name),

					resource.TestCheckResourceAttr("netbox_vlan.test", "vid", fmt.Sprintf("%d", vid)),
				),
			},
		},
	})

}

func TestAccVLANResource_full(t *testing.T) {

	name := testutil.RandomName("tf-test-vlan-full")

	vid := testutil.RandomVID()

	description := "Test VLAN with all fields"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterVLANCleanup(vid)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVLANDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccVLANResourceConfig_full(name, vid, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),

					resource.TestCheckResourceAttr("netbox_vlan.test", "name", name),

					resource.TestCheckResourceAttr("netbox_vlan.test", "vid", fmt.Sprintf("%d", vid)),

					resource.TestCheckResourceAttr("netbox_vlan.test", "description", description),

					resource.TestCheckResourceAttr("netbox_vlan.test", "status", "active"),
				),
			},
		},
	})

}

func TestAccVLANResource_withGroup(t *testing.T) {

	name := testutil.RandomName("tf-test-vlan-grp")

	vid := testutil.RandomVID()

	groupName := testutil.RandomName("tf-test-vlangrp")

	groupSlug := testutil.GenerateSlug("tf-test-vlangrp")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterVLANCleanup(vid)

	cleanup.RegisterVLANGroupCleanup(groupSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckVLANDestroy,

			testutil.CheckVLANGroupDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccVLANResourceConfig_withGroup(name, vid, groupName, groupSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),

					resource.TestCheckResourceAttr("netbox_vlan.test", "name", name),

					resource.TestCheckResourceAttr("netbox_vlan.test", "vid", fmt.Sprintf("%d", vid)),

					resource.TestCheckResourceAttrSet("netbox_vlan.test", "group"),
				),
			},
		},
	})

}

func TestAccVLANResource_update(t *testing.T) {

	name := testutil.RandomName("tf-test-vlan-upd")

	updatedName := testutil.RandomName("tf-test-vlan-updated")

	vid := testutil.RandomVID()

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterVLANCleanup(vid)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVLANDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccVLANResourceConfig_basic(name, vid),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),

					resource.TestCheckResourceAttr("netbox_vlan.test", "name", name),
				),
			},

			{

				Config: testAccVLANResourceConfig_full(updatedName, vid, "Updated description"),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),

					resource.TestCheckResourceAttr("netbox_vlan.test", "name", updatedName),

					resource.TestCheckResourceAttr("netbox_vlan.test", "description", "Updated description"),

					resource.TestCheckResourceAttr("netbox_vlan.test", "status", "active"),
				),
			},
		},
	})

}

func testAccVLANResourceConfig_basic(name string, vid int32) string {

	return fmt.Sprintf(`































resource "netbox_vlan" "test" {































  name = %q































  vid  = %d































}































`, name, vid)

}

func testAccVLANResourceConfig_full(name string, vid int32, description string) string {

	return fmt.Sprintf(`































resource "netbox_vlan" "test" {































  name        = %q































  vid         = %d































  description = %q































  status      = "active"































}































`, name, vid, description)

}

func testAccVLANResourceConfig_withGroup(name string, vid int32, groupName, groupSlug string) string {

	return fmt.Sprintf(`































resource "netbox_vlan_group" "test" {































  name = %q































  slug = %q































}































































resource "netbox_vlan" "test" {































  name  = %q































  vid   = %d































  group = netbox_vlan_group.test.id































}































`, groupName, groupSlug, name, vid)

}
