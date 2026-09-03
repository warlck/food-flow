package modifieroptionapi_test

import (
	"github.com/warlck/food-flow/app/domain/modifieroptionapp"
	"github.com/warlck/food-flow/business/domain/modifieroptionbus"
)

func toAppModifierOptionPtr(option modifieroptionbus.ModifierOption) *modifieroptionapp.ModifierOption {
	optionApp := modifieroptionapp.ToAppModifierOption(option)
	return &optionApp
}

func toAppModifierOption(option modifieroptionbus.ModifierOption) modifieroptionapp.ModifierOption {
	return modifieroptionapp.ToAppModifierOption(option)
}
