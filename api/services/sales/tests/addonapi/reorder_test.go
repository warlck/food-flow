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
			URL:        "/v1/addons/reorder",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusOK,
			// Swaps the middle pair. Addons[0] keeps rank 10 so update-200,
			// which runs later and asserts the seeded rank, is unaffected.
			Input: &addonapp.ReorderAddons{
				MenuItemID: sd.MenuItems[0].ID.String(),
				OrderedIDs: []string{
					sd.Addons[0].ID.String(),
					sd.Addons[2].ID.String(),
					sd.Addons[1].ID.String(),
					sd.Addons[3].ID.String(),
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
				if len(addons) != 4 {
					return fmt.Sprintf("expected 4 addons, got %d", len(addons))
				}

				expIDs := []string{
					sd.Addons[0].ID.String(),
					sd.Addons[2].ID.String(),
					sd.Addons[1].ID.String(),
					sd.Addons[3].ID.String(),
				}

				for i, expID := range expIDs {
					if addons[i].ID != expID {
						return fmt.Sprintf("position %d: expected addon %s, got %s", i, expID, addons[i].ID)
					}

					expRank := (i + 1) * 10
					if addons[i].Rank == nil || *addons[i].Rank != expRank {
						return fmt.Sprintf("position %d: rank is not %d", i, expRank)
					}
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
			Name:       "subset-mismatch",
			URL:        "/v1/addons/reorder",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &addonapp.ReorderAddons{
				MenuItemID: sd.MenuItems[0].ID.String(),
				OrderedIDs: []string{
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
		{
			Name:       "duplicate-id",
			URL:        "/v1/addons/reorder",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &addonapp.ReorderAddons{
				MenuItemID: sd.MenuItems[0].ID.String(),
				OrderedIDs: []string{
					sd.Addons[3].ID.String(),
					sd.Addons[2].ID.String(),
					sd.Addons[1].ID.String(),
					sd.Addons[1].ID.String(),
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
			Name:       "unknown-id",
			URL:        "/v1/addons/reorder",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &addonapp.ReorderAddons{
				MenuItemID: sd.MenuItems[0].ID.String(),
				OrderedIDs: []string{
					sd.Addons[3].ID.String(),
					sd.Addons[2].ID.String(),
					sd.Addons[1].ID.String(),
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
			Name:       "addon-from-other-item",
			URL:        "/v1/addons/reorder",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &addonapp.ReorderAddons{
				MenuItemID: sd.MenuItems[0].ID.String(),
				OrderedIDs: []string{
					sd.Addons[3].ID.String(),
					sd.Addons[2].ID.String(),
					sd.Addons[1].ID.String(),
					// Belongs to MenuItems[2]; same length but fails membership.
					sd.Addons[4].ID.String(),
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
			Name:       "unknown-menu-item",
			URL:        "/v1/addons/reorder",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &addonapp.ReorderAddons{
				MenuItemID: uuid.New().String(),
				OrderedIDs: []string{
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
			URL:        "/v1/addons/reorder",
			Token:      "",
			Method:     http.MethodPost,
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
			URL:        "/v1/addons/reorder",
			Token:      sd.Users[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusForbidden,
			Input: &addonapp.ReorderAddons{
				MenuItemID: sd.MenuItems[0].ID.String(),
				OrderedIDs: []string{
					sd.Addons[0].ID.String(),
					sd.Addons[1].ID.String(),
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
		{
			Name:       "other-org-item",
			URL:        "/v1/addons/reorder",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusForbidden,
			Input: &addonapp.ReorderAddons{
				MenuItemID: sd.MenuItems[2].ID.String(),
				OrderedIDs: []string{
					sd.Addons[4].ID.String(),
					sd.Addons[5].ID.String(),
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
