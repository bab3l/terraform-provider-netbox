// Package testutil provides utilities for acceptance testing of the Netbox provider.

package testutil

import (
	"context"
	"time"
)

// RegisterSiteGroupCleanup registers a cleanup function that will delete

// a site group by slug after the test completes.

func (c *CleanupResource) RegisterSiteGroupCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		// Find the site group by slug

		list, resp, err := c.client.DcimAPI.DcimSiteGroupsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			// Log but don't fail - resource might already be deleted

			c.t.Logf("Cleanup: failed to list site groups with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: site group with slug %s not found (already deleted)", slug)

			return

		}

		// Delete the site group

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimSiteGroupsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete site group %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted site group %d (slug: %s)", id, slug)

		}

	})

}

// RegisterSiteCleanup registers a cleanup function that will delete

// a site by slug after the test completes.

func (c *CleanupResource) RegisterSiteCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimSitesList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list sites with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: site with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimSitesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete site %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted site %d (slug: %s)", id, slug)

		}

	})

}

// RegisterManufacturerCleanup registers a cleanup function that will delete

// a manufacturer by slug after the test completes.

func (c *CleanupResource) RegisterManufacturerCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimManufacturersList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list manufacturers with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: manufacturer with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimManufacturersDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete manufacturer %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted manufacturer %d (slug: %s)", id, slug)

		}

	})

}

// RegisterPlatformCleanup registers a cleanup function that will delete

// a platform by slug after the test completes.

func (c *CleanupResource) RegisterPlatformCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimPlatformsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list platforms with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: platform with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimPlatformsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete platform %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted platform %d (slug: %s)", id, slug)

		}

	})

}

// RegisterRegionCleanup registers a cleanup function that will delete

// a region by slug after the test completes.

func (c *CleanupResource) RegisterRegionCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimRegionsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list regions with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: region with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimRegionsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete region %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted region %d (slug: %s)", id, slug)

		}

	})

}

// RegisterLocationCleanup registers a cleanup function that will delete

// a location by slug after the test completes.

func (c *CleanupResource) RegisterLocationCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimLocationsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list locations with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: location with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimLocationsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete location %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted location %d (slug: %s)", id, slug)

		}

	})

}

// RegisterRackCleanup registers a cleanup function that will delete

// a rack by name after the test completes.

func (c *CleanupResource) RegisterRackCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimRacksList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list racks with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: rack with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimRacksDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete rack %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted rack %d (name: %s)", id, name)

		}

	})

}

// RegisterDeviceRoleCleanup registers a cleanup function that will delete

// a device role by slug after the test completes.

func (c *CleanupResource) RegisterDeviceRoleCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimDeviceRolesList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list device roles with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: device role with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimDeviceRolesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete device role %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted device role %d (slug: %s)", id, slug)

		}

	})

}

// RegisterRackRoleCleanup registers a cleanup function that will delete

// a rack role by slug after the test completes.

func (c *CleanupResource) RegisterRackRoleCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimRackRolesList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list rack roles with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: rack role with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimRackRolesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete rack role %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted rack role %d (slug: %s)", id, slug)

		}

	})

}

// RegisterDeviceTypeCleanup registers a cleanup function that will delete

// a device type by slug after the test completes.

func (c *CleanupResource) RegisterDeviceTypeCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimDeviceTypesList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list device types with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: device type with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimDeviceTypesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete device type %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted device type %d (slug: %s)", id, slug)

		}

	})

}

// RegisterDeviceCleanup registers a cleanup function that will delete

// a device by name after the test completes.

func (c *CleanupResource) RegisterDeviceCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimDevicesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list devices with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: device with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimDevicesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete device %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted device %d (name: %s)", id, name)

		}

	})

}

// RegisterInterfaceCleanup registers a cleanup function that will delete

// an interface by name and device after the test completes.

func (c *CleanupResource) RegisterInterfaceCleanup(name string, deviceName string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		// First check if the parent device exists
		deviceList, deviceResp, deviceErr := c.client.DcimAPI.DcimDevicesList(ctx).Name([]string{deviceName}).Execute()
		if deviceErr != nil || deviceResp.StatusCode != 200 || len(deviceList.Results) == 0 {
			// Device doesn't exist, so interface is cascade-deleted
			c.t.Logf("Cleanup: interface with name %s not cleaned up (parent device %s already deleted)", name, deviceName)
			return
		}

		list, resp, err := c.client.DcimAPI.DcimInterfacesList(ctx).Name([]string{name}).Device([]*string{&deviceName}).Execute()

		if err != nil {
			// 404 means the interface doesn't exist
			if resp != nil && resp.StatusCode == 404 {
				c.t.Logf("Cleanup: interface with name %s on device %s not found (already deleted)", name, deviceName)
				return
			}

			c.t.Logf("Cleanup: failed to list interfaces with name %s on device %s: %v", name, deviceName, err)

			return

		}

		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: interface with name %s on device %s not found (already deleted)", name, deviceName)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimInterfacesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete interface %d (name: %s, device: %s): %v", id, name, deviceName, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted interface %d (name: %s, device: %s)", id, name, deviceName)

		}

	})

}

// RegisterRearPortTemplateCleanup registers a cleanup function that will delete

// a rear port template by name and device type after the test completes.

func (c *CleanupResource) RegisterRearPortTemplateCleanup(name string, deviceTypeID int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		dtID := deviceTypeID

		list, resp, err := c.client.DcimAPI.DcimRearPortTemplatesList(ctx).Name([]string{name}).DeviceTypeId([]*int32{&dtID}).Execute()

		if err != nil || resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: rear port template with name %s and device type %d not found: %v", name, deviceTypeID, err)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimRearPortTemplatesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete rear port template %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted rear port template %d (name: %s)", id, name)

		}

	})

}

// RegisterFrontPortTemplateCleanup registers a cleanup function that will delete

// a front port template by name and device type after the test completes.

func (c *CleanupResource) RegisterFrontPortTemplateCleanup(name string, deviceTypeID int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		dtID := deviceTypeID

		list, resp, err := c.client.DcimAPI.DcimFrontPortTemplatesList(ctx).Name([]string{name}).DeviceTypeId([]*int32{&dtID}).Execute()

		if err != nil || resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: front port template with name %s and device type %d not found: %v", name, deviceTypeID, err)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimFrontPortTemplatesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete front port template %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted front port template %d (name: %s)", id, name)

		}

	})

}

// RegisterRearPortCleanup registers a cleanup function that will delete

// a rear port by name and device after the test completes.

func (c *CleanupResource) RegisterRearPortCleanup(name string, deviceID int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimRearPortsList(ctx).Name([]string{name}).DeviceId([]int32{deviceID}).Execute()

		if err != nil || resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: rear port with name %s and device %d not found: %v", name, deviceID, err)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimRearPortsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete rear port %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted rear port %d (name: %s)", id, name)

		}

	})

}

// RegisterFrontPortCleanup registers a cleanup function that will delete

// a front port by name and device after the test completes.

func (c *CleanupResource) RegisterFrontPortCleanup(name string, deviceID int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimFrontPortsList(ctx).Name([]string{name}).DeviceId([]int32{deviceID}).Execute()

		if err != nil || resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: front port with name %s and device %d not found: %v", name, deviceID, err)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimFrontPortsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete front port %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted front port %d (name: %s)", id, name)

		}

	})

}

// RegisterDeviceBayTemplateCleanup registers a cleanup function that will delete

// a device bay template by name after the test completes.

func (c *CleanupResource) RegisterDeviceBayTemplateCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimDeviceBayTemplatesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list device bay templates with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: device bay template with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimDeviceBayTemplatesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete device bay template %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted device bay template %d (name: %s)", id, name)

		}

	})

}

// RegisterRackReservationCleanup registers a cleanup function that will delete

// a rack reservation by ID after the test completes.

func (c *CleanupResource) RegisterRackReservationCleanup(id int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		_, err := c.client.DcimAPI.DcimRackReservationsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete rack reservation %d: %v", id, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted rack reservation %d", id)

		}

	})

}

// RegisterVirtualDeviceContextCleanup registers a cleanup function that will delete

// a virtual device context by ID after the test completes.

func (c *CleanupResource) RegisterVirtualDeviceContextCleanup(id int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		_, err := c.client.DcimAPI.DcimVirtualDeviceContextsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete virtual device context %d: %v", id, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted virtual device context %d", id)

		}

	})

}

// RegisterModuleBayTemplateCleanup registers a cleanup function that will delete

// a module bay template by ID after the test completes.

func (c *CleanupResource) RegisterModuleBayTemplateCleanup(id int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		_, err := c.client.DcimAPI.DcimModuleBayTemplatesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete module bay template %d: %v", id, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted module bay template %d", id)

		}

	})

}

// RegisterInventoryItemTemplateCleanup registers a cleanup function that will delete

// an inventory item template by ID after the test completes.

func (c *CleanupResource) RegisterInventoryItemTemplateCleanup(id int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		_, err := c.client.DcimAPI.DcimInventoryItemTemplatesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete inventory item template %d: %v", id, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted inventory item template %d", id)

		}

	})

}

// RegisterPowerPanelCleanup registers a cleanup function that will delete

// a power panel by name after the test completes.

func (c *CleanupResource) RegisterPowerPanelCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimPowerPanelsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list power panels with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: power panel with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimPowerPanelsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete power panel %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted power panel %d (name: %s)", id, name)

		}

	})

}

// RegisterPowerFeedCleanup registers a cleanup function that will delete

// a power feed by name after the test completes.

func (c *CleanupResource) RegisterPowerFeedCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimPowerFeedsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list power feeds with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: power feed with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimPowerFeedsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete power feed %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted power feed %d (name: %s)", id, name)

		}

	})

}

// RegisterModuleBayCleanup registers a cleanup function that will delete

// a module bay by name and device ID after the test completes.

func (c *CleanupResource) RegisterModuleBayCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimModuleBaysList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list module bays with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: module bay with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimModuleBaysDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete module bay %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted module bay %d (name: %s)", id, name)

		}

	})

}

// RegisterVirtualChassisCleanup registers a cleanup function that will delete

// a virtual chassis by name after the test completes.

func (c *CleanupResource) RegisterVirtualChassisCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimVirtualChassisList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list virtual chassis with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: virtual chassis with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimVirtualChassisDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete virtual chassis %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted virtual chassis %d (name: %s)", id, name)

		}

	})

}

// RegisterInventoryItemCleanup registers an inventory item for cleanup during test execution.

func (c *CleanupResource) RegisterInventoryItemCleanup(id int32, name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		_, err := c.client.DcimAPI.DcimInventoryItemsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete inventory item %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted inventory item %d (name: %s)", id, name)

		}

	})

}

// RegisterInventoryItemRoleCleanup registers a cleanup function that will delete

// an inventory item role by name after the test completes.

func (c *CleanupResource) RegisterInventoryItemRoleCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimInventoryItemRolesList(ctx).NameIc([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list inventory item roles with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: inventory item role with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimInventoryItemRolesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete inventory item role %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted inventory item role %d (name: %s)", id, name)

		}

	})

}
