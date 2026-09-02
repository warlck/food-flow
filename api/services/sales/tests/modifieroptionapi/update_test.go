package modifieroptionapi_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/domain/modifieroptionapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
)

func update404(sd apitest.SeedData) []apitest.Table {
	unknownID := uuid.New()

	table := []apitest.Table{
		{
			Name:       "not-found",
			URL:        fmt.Sprintf("/v1/modifieroptions/%s", unknownID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusNotFound,
			Input:      &modifieroptionapp.UpdateModifierOption{},
			GotResp:    &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.NotFound,
				Message: fmt.Sprintf("query: optionID[%s]: db: modifier option not found", unknownID),
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
