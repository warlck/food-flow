package addonapi_test

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/domain/addonapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
)

func reorder200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/addons/order",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &addonapp.ReorderAddons{
				CategoryID: sd.Categories[1].ID.String(),
				OrderedIDs: []string{
					sd.Addons[3].ID.String(),
					sd.Addons[2].ID.String(),
				},
			},
			GotResp: &[]addonapp.Addon{},
			ExpResp: &[]addonapp.Addon{},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*[]addonapp.Addon)
				if !exists {
					return "got is not *[]addonapp.Addon"
				}

				addons := *gotResp
				if len(addons) != 2 {
					return fmt.Sprintf("expected 2 addons, got %d", len(addons))
				}

				if addons[0].ID != sd.Addons[3].ID.String() || addons[1].ID != sd.Addons[2].ID.String() {
					return fmt.Sprintf("unexpected order: %s, %s", addons[0].ID, addons[1].ID)
				}

				if addons[0].Rank == nil || *addons[0].Rank != 10 {
					return "first addon rank is not 10"
				}

				if addons[1].Rank == nil || *addons[1].Rank != 20 {
					return "second addon rank is not 20"
				}

				return ""
			},
		},
	}

	return table
}

func reorder400(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "mismatch-length",
			URL:        "/v1/addons/order",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &addonapp.ReorderAddons{
				CategoryID: sd.Categories[1].ID.String(),
				OrderedIDs: []string{
					sd.Addons[2].ID.String(),
				},
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code: errs.InvalidArgument,
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(*errs.Error)
				if gotErr.Code != errs.InvalidArgument {
					return "error code mismatch"
				}
				return ""
			},
		},
		{
			Name:       "invalid-id",
			URL:        "/v1/addons/order",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &addonapp.ReorderAddons{
				CategoryID: sd.Categories[1].ID.String(),
				OrderedIDs: []string{
					sd.Addons[2].ID.String(),
					uuid.New().String(),
				},
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code: errs.InvalidArgument,
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(*errs.Error)
				if gotErr.Code != errs.InvalidArgument {
					return "error code mismatch"
				}
				return ""
			},
		},
		{
			Name:       "wrong-category-id",
			URL:        "/v1/addons/order",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &addonapp.ReorderAddons{
				CategoryID: sd.Categories[1].ID.String(),
				OrderedIDs: []string{
					sd.Addons[3].ID.String(),
					// Belongs to Categories[0]; same length but fails membership.
					sd.Addons[0].ID.String(),
				},
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code: errs.InvalidArgument,
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(*errs.Error)
				if gotErr.Code != errs.InvalidArgument {
					return "error code mismatch"
				}
				return ""
			},
		},
	}

	return table
}

func reorder401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "emptytoken",
			URL:        "/v1/addons/order",
			Token:      "",
			Method:     http.MethodPut,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.Unauthenticated,
				Message: "expected authorization header format: Bearer <token>",
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(*errs.Error)
				expErr := exp.(*errs.Error)
				if gotErr.Code != expErr.Code {
					return "error code mismatch"
				}
				if gotErr.Message != expErr.Message {
					return "error message mismatch"
				}
				return ""
			},
		},
		{
			Name:       "wronguser",
			URL:        "/v1/addons/order",
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusForbidden,
			Input: &addonapp.ReorderAddons{
				CategoryID: sd.Categories[1].ID.String(),
				OrderedIDs: []string{
					sd.Addons[2].ID.String(),
					sd.Addons[3].ID.String(),
				},
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.PermissionDenied,
				Message: "authorize: you are not authorized for that action, claims[[USER]] rule[rule_admin_only]",
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(*errs.Error)
				expErr := exp.(*errs.Error)
				if gotErr.Code != expErr.Code {
					return "error code mismatch"
				}
				if gotErr.Message != expErr.Message {
					return "error message mismatch"
				}
				return ""
			},
		},
	}

	return table
}
