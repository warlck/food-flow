package modifiergroupapi_test

import (
	"github.com/warlck/food-flow/app/domain/modifiergroupapp"
	"github.com/warlck/food-flow/business/domain/modifiergroupbus"
)

func toAppModifierGroupPtr(group modifiergroupbus.ModifierGroup) *modifiergroupapp.ModifierGroup {
	groupApp := modifiergroupapp.ToAppModifierGroup(group)
	return &groupApp
}

func toAppModifierGroup(group modifiergroupbus.ModifierGroup) modifiergroupapp.ModifierGroup {
	return modifiergroupapp.ToAppModifierGroup(group)
}
