package organizationapi_test

import (
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/app/domain/organizationapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
)

func queryMyOrgs200(sd apitest.SeedData) []apitest.Table {
	orgs := make([]organizationapp.Organization, 0, len(sd.Organizations))
	for _, org := range sd.Organizations {
		orgs = append(orgs, organizationapp.ToAppOrganization(org.Organization))
	}

	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/organizations/me",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &[]organizationapp.Organization{},
			ExpResp:    &orgs,
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*[]organizationapp.Organization)
				expResp := exp.(*[]organizationapp.Organization)
				if len(*gotResp) == len(*expResp) {
					for i := range *gotResp {
						(*expResp)[i].DateCreated = (*gotResp)[i].DateCreated
						(*expResp)[i].DateUpdated = (*gotResp)[i].DateUpdated
					}
				}
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func queryMyOrgs401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "emptytoken",
			URL:        "/v1/organizations/me",
			Token:      "&nbsp;",
			StatusCode: http.StatusUnauthorized,
			Method:     http.MethodGet,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "error parsing token: token is malformed: token contains an invalid number of segments"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "wronguser",
			URL:        "/v1/organizations/me",
			Token:      sd.Users[0].Token,
			StatusCode: http.StatusForbidden,
			Method:     http.MethodGet,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.PermissionDenied, "authorize: you are not authorized for that action, claims[[USER]] rule[rule_admin_only]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
