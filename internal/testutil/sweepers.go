package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const sweeperPrefix = "tf-test"

func isSweeperMatch(value string) bool {

	v := strings.ToLower(strings.TrimSpace(value))

	return v != "" && strings.HasPrefix(v, sweeperPrefix)

}

func init() {

	resource.AddTestSweepers("netbox_device", &resource.Sweeper{

		Name: "netbox_device",

		F: sweepDcimDevices,
	})

	resource.AddTestSweepers("netbox_device_type", &resource.Sweeper{

		Name: "netbox_device_type",

		Dependencies: []string{"netbox_device"},

		F: sweepDcimDeviceTypes,
	})

	resource.AddTestSweepers("netbox_device_role", &resource.Sweeper{

		Name: "netbox_device_role",

		Dependencies: []string{"netbox_device"},

		F: sweepDcimDeviceRoles,
	})

	resource.AddTestSweepers("netbox_site", &resource.Sweeper{

		Name: "netbox_site",

		Dependencies: []string{"netbox_device"},

		F: sweepDcimSites,
	})

	resource.AddTestSweepers("netbox_site_group", &resource.Sweeper{

		Name: "netbox_site_group",

		Dependencies: []string{"netbox_site"},

		F: sweepDcimSiteGroups,
	})

	resource.AddTestSweepers("netbox_tenant", &resource.Sweeper{

		Name: "netbox_tenant",

		F: sweepTenancyTenants,
	})

	resource.AddTestSweepers("netbox_tenant_group", &resource.Sweeper{

		Name: "netbox_tenant_group",

		Dependencies: []string{"netbox_tenant"},

		F: sweepTenancyTenantGroups,
	})

	resource.AddTestSweepers("netbox_manufacturer", &resource.Sweeper{

		Name: "netbox_manufacturer",

		Dependencies: []string{"netbox_device_type"},

		F: sweepDcimManufacturers,
	})

}

func sweepDcimSiteGroups(_ string) error {

	client, err := GetSharedClient()

	if err != nil {

		return err

	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

	defer cancel()

	var (
		limit int32 = 200

		offset int32
	)

	for {

		list, _, err := client.DcimAPI.DcimSiteGroupsList(ctx).
			Limit(limit).
			Offset(offset).
			SlugIsw([]string{sweeperPrefix}).
			Execute()

		if err != nil {

			return fmt.Errorf("list site groups: %w", err)

		}

		for _, sg := range list.Results {

			slug := sg.GetSlug()

			name := sg.GetName()

			if !isSweeperMatch(slug) && !isSweeperMatch(name) {

				continue

			}

			if _, err := client.DcimAPI.DcimSiteGroupsDestroy(ctx, sg.GetId()).Execute(); err != nil {

				return fmt.Errorf("delete site group %q (id=%d): %w", slug, sg.GetId(), err)

			}

		}

		offset += int32(len(list.Results)) // #nosec G115 -- len(list.Results) is bounded by Limit(200); safe conversion for pagination offset.

		if offset >= list.GetCount() || len(list.Results) == 0 {

			break

		}

	}

	return nil

}

func sweepDcimSites(_ string) error {

	client, err := GetSharedClient()

	if err != nil {

		return err

	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

	defer cancel()

	var (
		limit int32 = 200

		offset int32
	)

	for {

		list, _, err := client.DcimAPI.DcimSitesList(ctx).
			Limit(limit).
			Offset(offset).
			SlugIsw([]string{sweeperPrefix}).
			Execute()

		if err != nil {

			return fmt.Errorf("list sites: %w", err)

		}

		for _, s := range list.Results {

			slug := s.GetSlug()

			name := s.GetName()

			if !isSweeperMatch(slug) && !isSweeperMatch(name) {

				continue

			}

			if _, err := client.DcimAPI.DcimSitesDestroy(ctx, s.GetId()).Execute(); err != nil {

				return fmt.Errorf("delete site %q (id=%d): %w", slug, s.GetId(), err)

			}

		}

		offset += int32(len(list.Results)) // #nosec G115 -- len(list.Results) is bounded by Limit(200); safe conversion for pagination offset.

		if offset >= list.GetCount() || len(list.Results) == 0 {

			break

		}

	}

	return nil

}

func sweepDcimManufacturers(_ string) error {

	client, err := GetSharedClient()

	if err != nil {

		return err

	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

	defer cancel()

	var (
		limit int32 = 200

		offset int32
	)

	for {

		list, _, err := client.DcimAPI.DcimManufacturersList(ctx).
			Limit(limit).
			Offset(offset).
			SlugIsw([]string{sweeperPrefix}).
			Execute()

		if err != nil {

			return fmt.Errorf("list manufacturers: %w", err)

		}

		for _, m := range list.Results {

			slug := m.GetSlug()

			name := m.GetName()

			if !isSweeperMatch(slug) && !isSweeperMatch(name) {

				continue

			}

			if _, err := client.DcimAPI.DcimManufacturersDestroy(ctx, m.GetId()).Execute(); err != nil {

				return fmt.Errorf("delete manufacturer %q (id=%d): %w", slug, m.GetId(), err)

			}

		}

		offset += int32(len(list.Results)) // #nosec G115 -- len(list.Results) is bounded by Limit(200); safe conversion for pagination offset.

		if offset >= list.GetCount() || len(list.Results) == 0 {

			break

		}

	}

	return nil

}

func sweepDcimDeviceRoles(_ string) error {

	client, err := GetSharedClient()

	if err != nil {

		return err

	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

	defer cancel()

	var (
		limit int32 = 200

		offset int32
	)

	for {

		list, _, err := client.DcimAPI.DcimDeviceRolesList(ctx).
			Limit(limit).
			Offset(offset).
			SlugIsw([]string{sweeperPrefix}).
			Execute()

		if err != nil {

			return fmt.Errorf("list device roles: %w", err)

		}

		for _, r := range list.Results {

			slug := r.GetSlug()

			name := r.GetName()

			if !isSweeperMatch(slug) && !isSweeperMatch(name) {

				continue

			}

			if _, err := client.DcimAPI.DcimDeviceRolesDestroy(ctx, r.GetId()).Execute(); err != nil {

				return fmt.Errorf("delete device role %q (id=%d): %w", slug, r.GetId(), err)

			}

		}

		offset += int32(len(list.Results)) // #nosec G115 -- len(list.Results) is bounded by Limit(200); safe conversion for pagination offset.

		if offset >= list.GetCount() || len(list.Results) == 0 {

			break

		}

	}

	return nil

}

func sweepDcimDeviceTypes(_ string) error {

	client, err := GetSharedClient()

	if err != nil {

		return err

	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

	defer cancel()

	var (
		limit int32 = 200

		offset int32
	)

	for {

		list, _, err := client.DcimAPI.DcimDeviceTypesList(ctx).
			Limit(limit).
			Offset(offset).
			SlugIsw([]string{sweeperPrefix}).
			Execute()

		if err != nil {

			return fmt.Errorf("list device types: %w", err)

		}

		for _, dt := range list.Results {

			slug := dt.GetSlug()

			model := dt.GetModel()

			if !isSweeperMatch(slug) && !isSweeperMatch(model) {

				continue

			}

			if _, err := client.DcimAPI.DcimDeviceTypesDestroy(ctx, dt.GetId()).Execute(); err != nil {

				return fmt.Errorf("delete device type %q (id=%d): %w", slug, dt.GetId(), err)

			}

		}

		offset += int32(len(list.Results)) // #nosec G115 -- len(list.Results) is bounded by Limit(200); safe conversion for pagination offset.

		if offset >= list.GetCount() || len(list.Results) == 0 {

			break

		}

	}

	return nil

}

func sweepDcimDevices(_ string) error {

	client, err := GetSharedClient()

	if err != nil {

		return err

	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

	defer cancel()

	var (
		limit int32 = 200

		offset int32
	)

	for {

		list, _, err := client.DcimAPI.DcimDevicesList(ctx).
			Limit(limit).
			Offset(offset).
			NameIsw([]string{sweeperPrefix}).
			Execute()

		if err != nil {

			return fmt.Errorf("list devices: %w", err)

		}

		for _, d := range list.Results {

			name := d.GetName()

			if !isSweeperMatch(name) {

				continue

			}

			if _, err := client.DcimAPI.DcimDevicesDestroy(ctx, d.GetId()).Execute(); err != nil {

				return fmt.Errorf("delete device %q (id=%d): %w", name, d.GetId(), err)

			}

		}

		offset += int32(len(list.Results)) // #nosec G115 -- len(list.Results) is bounded by Limit(200); safe conversion for pagination offset.

		if offset >= list.GetCount() || len(list.Results) == 0 {

			break

		}

	}

	return nil

}

func sweepTenancyTenants(_ string) error {

	client, err := GetSharedClient()

	if err != nil {

		return err

	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

	defer cancel()

	var (
		limit int32 = 200

		offset int32
	)

	for {

		list, _, err := client.TenancyAPI.TenancyTenantsList(ctx).
			Limit(limit).
			Offset(offset).
			SlugIsw([]string{sweeperPrefix}).
			Execute()

		if err != nil {

			return fmt.Errorf("list tenants: %w", err)

		}

		for _, t := range list.Results {

			slug := t.GetSlug()

			name := t.GetName()

			if !isSweeperMatch(slug) && !isSweeperMatch(name) {

				continue

			}

			if _, err := client.TenancyAPI.TenancyTenantsDestroy(ctx, t.GetId()).Execute(); err != nil {

				return fmt.Errorf("delete tenant %q (id=%d): %w", slug, t.GetId(), err)

			}

		}

		offset += int32(len(list.Results)) // #nosec G115 -- len(list.Results) is bounded by Limit(200); safe conversion for pagination offset.

		if offset >= list.GetCount() || len(list.Results) == 0 {

			break

		}

	}

	return nil

}

func sweepTenancyTenantGroups(_ string) error {

	client, err := GetSharedClient()

	if err != nil {

		return err

	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

	defer cancel()

	var (
		limit int32 = 200

		offset int32
	)

	for {

		list, _, err := client.TenancyAPI.TenancyTenantGroupsList(ctx).
			Limit(limit).
			Offset(offset).
			SlugIsw([]string{sweeperPrefix}).
			Execute()

		if err != nil {

			return fmt.Errorf("list tenant groups: %w", err)

		}

		for _, tg := range list.Results {

			slug := tg.GetSlug()

			name := tg.GetName()

			if !isSweeperMatch(slug) && !isSweeperMatch(name) {

				continue

			}

			if _, err := client.TenancyAPI.TenancyTenantGroupsDestroy(ctx, tg.GetId()).Execute(); err != nil {

				return fmt.Errorf("delete tenant group %q (id=%d): %w", slug, tg.GetId(), err)

			}

		}

		offset += int32(len(list.Results)) // #nosec G115 -- len(list.Results) is bounded by Limit(200); safe conversion for pagination offset.

		if offset >= list.GetCount() || len(list.Results) == 0 {

			break

		}

	}

	return nil

}
