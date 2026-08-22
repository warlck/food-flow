package menuitemapi_test

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/domain/menuitemapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
)

func reorder200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/menuitems/order",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &menuitemapp.ReorderMenuItems{
				CategoryID: sd.Categories[1].ID.String(),
				OrderedIDs: []string{
					sd.MenuItems[3].ID.String(),
					sd.MenuItems[2].ID.String(),
				},
			},
			GotResp: &[]menuitemapp.MenuItem{},
			ExpResp: &[]menuitemapp.MenuItem{},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*[]menuitemapp.MenuItem)
				if !exists {
					return "got is not *[]menuitemapp.MenuItem"
				}

				items := *gotResp
				if len(items) != 2 {
					return fmt.Sprintf("expected 2 items, got %d", len(items))
				}

				if items[0].ID != sd.MenuItems[3].ID.String() || items[1].ID != sd.MenuItems[2].ID.String() {
					return fmt.Sprintf("unexpected order: %s, %s", items[0].ID, items[1].ID)
				}

				if items[0].Rank == nil || *items[0].Rank != 10 {
					return "first item rank is not 10"
				}

				if items[1].Rank == nil || *items[1].Rank != 20 {
					return "second item rank is not 20"
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
			URL:        "/v1/menuitems/order",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &menuitemapp.ReorderMenuItems{
				CategoryID: sd.Categories[1].ID.String(),
				OrderedIDs: []string{
					sd.MenuItems[2].ID.String(),
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
			URL:        "/v1/menuitems/order",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &menuitemapp.ReorderMenuItems{
				CategoryID: sd.Categories[1].ID.String(),
				OrderedIDs: []string{
					sd.MenuItems[2].ID.String(),
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
			URL:        "/v1/menuitems/order",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &menuitemapp.ReorderMenuItems{
				CategoryID: sd.Categories[1].ID.String(),
				OrderedIDs: []string{
					sd.MenuItems[3].ID.String(),
					// Belongs to Categories[0]; same length but fails membership.
					sd.MenuItems[0].ID.String(),
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
			URL:        "/v1/menuitems/order",
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
			URL:        "/v1/menuitems/order",
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusForbidden,
			Input: &menuitemapp.ReorderMenuItems{
				CategoryID: sd.Categories[1].ID.String(),
				OrderedIDs: []string{
					sd.MenuItems[2].ID.String(),
					sd.MenuItems[3].ID.String(),
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
		{
			Name:       "other-org-category",
			URL:        "/v1/menuitems/order",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusForbidden,
			Input: &menuitemapp.ReorderMenuItems{
				CategoryID: sd.Categories[3].ID.String(),
				OrderedIDs: []string{
					sd.MenuItems[6].ID.String(),
					sd.MenuItems[7].ID.String(),
				},
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.PermissionDenied,
				Message: fmt.Sprintf("user not in organization %s", sd.Organizations[1].ID),
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
