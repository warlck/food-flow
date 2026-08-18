package imageapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/domain/imageapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/query"
)

func Test_ImageAPI(t *testing.T) {
	t.Parallel()

	test := apitest.New(t, "Test_ImageAPI")

	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	test.Run(t, uploadUrlTable(sd), "uploadurl")
	test.Run(t, completeDeleteTable(sd), "complete-delete")
	t.Run("flow", func(t *testing.T) { uploadFlow(t, test, sd) })
}

// codeOnlyCmp compares only the error code so tests do not depend on
// validator-generated message wording.
func codeOnlyCmp(got any, exp any) string {
	gotResp, ok := got.(*errs.Error)
	if !ok {
		return "got is not an *errs.Error"
	}
	expResp, ok := exp.(*errs.Error)
	if !ok {
		return "exp is not an *errs.Error"
	}
	if gotResp.Code != expResp.Code {
		return "code mismatch: got " + gotResp.Code.String() + " want " + expResp.Code.String()
	}
	return ""
}

func uploadUrlTable(sd apitest.SeedData) []apitest.Table {
	restID := sd.Restaurants[0].ID

	table := []apitest.Table{
		{
			Name:       "missing-restaurant-id",
			URL:        "/v1/images/upload-url",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &imageapp.NewUpload{
				EntityType:  "menu_item",
				ContentType: "image/jpeg",
				SizeBytes:   100,
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, ""),
			CmpFunc: codeOnlyCmp,
		},
		{
			Name:       "invalid-entity-type",
			URL:        "/v1/images/upload-url",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &imageapp.NewUpload{
				RestaurantID: restID,
				EntityType:   "addon",
				ContentType:  "image/jpeg",
				SizeBytes:    100,
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, ""),
			CmpFunc: codeOnlyCmp,
		},
		{
			Name:       "invalid-content-type",
			URL:        "/v1/images/upload-url",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &imageapp.NewUpload{
				RestaurantID: restID,
				EntityType:   "menu_item",
				ContentType:  "image/gif",
				SizeBytes:    100,
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, ""),
			CmpFunc: codeOnlyCmp,
		},
		{
			Name:       "zero-size",
			URL:        "/v1/images/upload-url",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &imageapp.NewUpload{
				RestaurantID: restID,
				EntityType:   "menu_item",
				ContentType:  "image/png",
				SizeBytes:    0,
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, ""),
			CmpFunc: codeOnlyCmp,
		},
		{
			Name:       "size-over-limit",
			URL:        "/v1/images/upload-url",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &imageapp.NewUpload{
				RestaurantID: restID,
				EntityType:   "menu_item",
				ContentType:  "image/png",
				SizeBytes:    5242881,
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "invalid image data: size_bytes 5242881 exceeds the limit of 5242880 bytes"),
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*errs.Error)
				expResp := exp.(*errs.Error)
				if !gotResp.Equal(expResp) {
					return "got " + gotResp.Message + " want " + expResp.Message
				}
				return ""
			},
		},
		{
			Name:       "bad-restaurant-uuid",
			URL:        "/v1/images/upload-url",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: map[string]any{
				"restaurantId": "not-a-uuid",
				"entityType":   "restaurant",
				"contentType":  "image/jpeg",
				"sizeBytes":    100,
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, ""),
			CmpFunc: codeOnlyCmp,
		},
		{
			Name:       "non-admin-forbidden",
			URL:        "/v1/images/upload-url",
			Token:      sd.Users[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			Input: &imageapp.NewUpload{
				RestaurantID: restID,
				EntityType:   "restaurant",
				ContentType:  "image/jpeg",
				SizeBytes:    100,
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.Unauthenticated, ""),
			CmpFunc: codeOnlyCmp,
		},
	}

	return table
}

func completeDeleteTable(sd apitest.SeedData) []apitest.Table {
	missingID := uuid.NewString()

	table := []apitest.Table{
		{
			Name:       "complete-not-found",
			URL:        "/v1/images/" + missingID + "/complete",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusNotFound,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.NotFound, "query by id: image not found"),
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*errs.Error)
				expResp := exp.(*errs.Error)
				if !gotResp.Equal(expResp) {
					return "got " + gotResp.Message + " want " + expResp.Message
				}
				return ""
			},
		},
		{
			Name:       "delete-not-found",
			URL:        "/v1/images/" + missingID,
			Token:      sd.Admins[0].Token,
			Method:     http.MethodDelete,
			StatusCode: http.StatusNotFound,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.NotFound, "query by id: image not found"),
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*errs.Error)
				expResp := exp.(*errs.Error)
				if !gotResp.Equal(expResp) {
					return "got " + gotResp.Message + " want " + expResp.Message
				}
				return ""
			},
		},
	}

	return table
}

// uploadFlow exercises the full local-backend upload lifecycle end to end:
// mint a signed URL, PUT the bytes, confirm, list, download, and delete.
func uploadFlow(t *testing.T, test *apitest.Test, sd apitest.SeedData) {
	h := test.Handler()
	token := sd.Admins[0].Token
	restID := sd.Restaurants[0].ID

	do := func(method, url, contentType string, body []byte, decode any, wantStatus int) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(method, url, bytes.NewReader(body))
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != wantStatus {
			t.Fatalf("%s %s: got status %d, want %d; body: %s", method, url, w.Code, wantStatus, w.Body.String())
		}
		if decode != nil {
			if err := json.Unmarshal(w.Body.Bytes(), decode); err != nil {
				t.Fatalf("%s %s: decoding response: %v", method, url, err)
			}
		}
		return w
	}

	// 1. Mint the signed upload URL.
	newUpload := imageapp.NewUpload{
		RestaurantID: restID,
		EntityType:   "menu_item",
		ContentType:  "image/png",
		SizeBytes:    11,
	}
	grantBody, err := json.Marshal(newUpload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var grant imageapp.UploadGrant
	do(http.MethodPost, "/v1/images/upload-url", "application/json", grantBody, &grant, http.StatusCreated)

	if !strings.HasPrefix(grant.UploadURL, "/v1/images/local/menu-items%2F"+restID.String()) {
		t.Fatalf("unexpected upload url: %s", grant.UploadURL)
	}
	if !strings.HasSuffix(grant.UploadURL, ".png") {
		t.Fatalf("upload url should carry the png extension: %s", grant.UploadURL)
	}
	if grant.Image.Status != "pending" {
		t.Fatalf("image status = %q, want pending", grant.Image.Status)
	}

	// 2. Confirming before the bytes land fails with 400.
	do(http.MethodPost, "/v1/images/"+grant.Image.ID+"/complete", "", nil, nil, http.StatusBadRequest)

	// 3. PUT the bytes to the signed URL (local backend stores them).
	payload := []byte("hello-image")
	do(http.MethodPut, grant.UploadURL, "image/png", payload, nil, http.StatusOK)

	// 4. Confirm flips status and adopts the actual object size.
	var confirmed imageapp.Image
	do(http.MethodPost, "/v1/images/"+grant.Image.ID+"/complete", "", nil, &confirmed, http.StatusOK)
	if confirmed.Status != "confirmed" {
		t.Fatalf("status = %q, want confirmed", confirmed.Status)
	}
	if confirmed.SizeBytes != int64(len(payload)) {
		t.Fatalf("size = %d, want %d", confirmed.SizeBytes, len(payload))
	}
	if confirmed.PublicURL != grant.UploadURL {
		t.Fatalf("public url = %q, want %q", confirmed.PublicURL, grant.UploadURL)
	}

	// 5. Query lists the confirmed image for the restaurant.
	var list query.Result[imageapp.Image]
	do(http.MethodGet, "/v1/images?restaurant_id="+restID.String()+"&status=confirmed", "", nil, &list, http.StatusOK)
	if list.Total != 1 || len(list.Items) != 1 || list.Items[0].ID != confirmed.ID {
		t.Fatalf("list = %+v, want the confirmed image", list)
	}

	// 6. Download returns the uploaded bytes.
	dl := do(http.MethodGet, confirmed.PublicURL, "", nil, nil, http.StatusOK)
	if !bytes.Equal(dl.Body.Bytes(), payload) {
		t.Fatalf("download body = %q, want %q", dl.Body.String(), payload)
	}

	// 7. Delete removes the row and the object.
	do(http.MethodDelete, "/v1/images/"+confirmed.ID, "", nil, nil, http.StatusNoContent)
	do(http.MethodGet, confirmed.PublicURL, "", nil, nil, http.StatusNotFound)

	var after query.Result[imageapp.Image]
	do(http.MethodGet, "/v1/images?restaurant_id="+restID.String(), "", nil, &after, http.StatusOK)
	if after.Total != 0 {
		t.Fatalf("list after delete: total = %d, want 0", after.Total)
	}
}
