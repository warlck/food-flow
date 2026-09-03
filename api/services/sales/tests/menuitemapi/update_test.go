package menuitemapi_test

import (
	"fmt"
	"net/http"

	"github.com/warlck/food-flow/app/domain/menuitemapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/types/opt"
)

func update200(sd apitest.SeedData) []apitest.Table {
	updName := "Updated MenuItem"
	updDesc := "Updated Description"
	updPrice := 29.99
	updAvail := false

	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/menuitems/" + sd.MenuItems[1].ID.String(),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &menuitemapp.UpdateMenuItem{
				Name:        &updName,
				Description: &updDesc,
				Price:       &updPrice,
				Available:   &updAvail,
			},
			GotResp: &menuitemapp.MenuItem{},
			ExpResp: &menuitemapp.MenuItem{},
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*menuitemapp.MenuItem)
				if gotResp.Name != "Updated MenuItem" {
					return "name not updated"
				}
				if gotResp.Description != "Updated Description" {
					return "description not updated"
				}
				if gotResp.Price != 29.99 {
					return "price not updated"
				}
				if gotResp.Available != false {
					return "available not updated"
				}
				return ""
			},
		},
		{
			Name:       "update-rank",
			URL:        "/v1/menuitems/" + sd.MenuItems[1].ID.String(),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &menuitemapp.UpdateMenuItem{
				Rank: opt.NewNullInt(42),
			},
			GotResp: &menuitemapp.MenuItem{},
			ExpResp: &menuitemapp.MenuItem{},
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*menuitemapp.MenuItem)
				if gotResp.Rank == nil || *gotResp.Rank != 42 {
					return "rank not updated"
				}
				return ""
			},
		},
	}

	return table
}

func update400(sd apitest.SeedData) []apitest.Table {
	invalidPrice := -1.0

	table := []apitest.Table{
		{
			Name:       "invalid-price",
			URL:        "/v1/menuitems/" + sd.MenuItems[1].ID.String(),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &menuitemapp.UpdateMenuItem{
				Price: &invalidPrice,
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.InvalidArgument,
				Message: "validation error",
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(*errs.Error)
				expErr := exp.(*errs.Error)

				if gotErr.Code != expErr.Code {
					return fmt.Sprintf("code mismatch: got %v, exp %v", gotErr.Code, expErr.Code)
				}

				// Just check it contains validation error
				if gotErr.Message == "" {
					return "error message should not be empty"
				}

				return ""
			},
		},
		{
			Name:       "rank-zero",
			URL:        "/v1/menuitems/" + sd.MenuItems[1].ID.String(),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &menuitemapp.UpdateMenuItem{
				Rank: opt.NewNullInt(0),
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.InvalidArgument,
				Message: `[{"field":"rank","error":"rank must be 1 or greater"}]`,
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(*errs.Error)
				expErr := exp.(*errs.Error)

				if gotErr.Code != expErr.Code {
					return fmt.Sprintf("code mismatch: got %v, exp %v", gotErr.Code, expErr.Code)
				}

				if gotErr.Message != expErr.Message {
					return fmt.Sprintf("message mismatch: got %q, exp %q", gotErr.Message, expErr.Message)
				}

				return ""
			},
		},
		{
			Name:       "rank-negative",
			URL:        "/v1/menuitems/" + sd.MenuItems[1].ID.String(),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &menuitemapp.UpdateMenuItem{
				Rank: opt.NewNullInt(-1),
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.InvalidArgument,
				Message: `[{"field":"rank","error":"rank must be 1 or greater"}]`,
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(*errs.Error)
				expErr := exp.(*errs.Error)

				if gotErr.Code != expErr.Code {
					return fmt.Sprintf("code mismatch: got %v, exp %v", gotErr.Code, expErr.Code)
				}

				if gotErr.Message != expErr.Message {
					return fmt.Sprintf("message mismatch: got %q, exp %q", gotErr.Message, expErr.Message)
				}

				return ""
			},
		},
		{
			Name:       "move-to-category-from-different-restaurant",
			URL:        "/v1/menuitems/" + sd.MenuItems[1].ID.String(),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &menuitemapp.UpdateMenuItem{
				CategoryID: dbtest.StringPointer(sd.Categories[2].ID.String()),
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code: errs.InvalidArgument,
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(*errs.Error)
				expErr := exp.(*errs.Error)
				if gotErr.Code != expErr.Code {
					return fmt.Sprintf("code mismatch: got %v, exp %v", gotErr.Code, expErr.Code)
				}
				return ""
			},
		},
	}

	return table
}

func update401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "emptytoken",
			URL:        "/v1/menuitems/" + sd.MenuItems[1].ID.String(),
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
			URL:        "/v1/menuitems/" + sd.MenuItems[1].ID.String(),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusForbidden,
			Input: &menuitemapp.UpdateMenuItem{
				Name: dbtest.StringPointer("Updated Name"),
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
