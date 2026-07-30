package addonapi_test

import (
	"github.com/warlck/food-flow/app/domain/addonapp"
	"github.com/warlck/food-flow/business/domain/addonbus"
)

func toAppAddonPtr(addon addonbus.Addon) *addonapp.Addon {
	addonApp := addonapp.ToAppAddon(addon)
	return &addonApp
}

func toAppAddon(addon addonbus.Addon) addonapp.Addon {
	return addonapp.ToAppAddon(addon)
}
