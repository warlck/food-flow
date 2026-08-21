package imagebus_test

import (
	"context"
	"errors"
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/imagebus"
	"github.com/warlck/food-flow/business/domain/imagebus/stores/imagedb"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/storage"
)

// fakeSigner is an in-memory storage.Signer for tests.
type fakeSigner struct {
	objects map[string]int64
}

func newFakeSigner() *fakeSigner {
	return &fakeSigner{objects: make(map[string]int64)}
}

func (f *fakeSigner) SignUpload(_ context.Context, req storage.SignRequest) (storage.SignedUpload, error) {
	return storage.SignedUpload{
		UploadURL:  "https://storage.example.test/upload/" + req.ObjectPath,
		Method:     "PUT",
		Headers:    map[string]string{"Content-Type": req.ContentType},
		PublicURL:  "https://cdn.example.test/" + req.ObjectPath,
		ObjectPath: req.ObjectPath,
		ExpiresAt:  time.Now().Add(15 * time.Minute),
	}, nil
}

func (f *fakeSigner) Stat(_ context.Context, objectPath string) (storage.ObjectInfo, error) {
	size, ok := f.objects[objectPath]
	if !ok {
		return storage.ObjectInfo{}, storage.ErrNotFound
	}
	return storage.ObjectInfo{ObjectPath: objectPath, SizeBytes: size, ContentType: "image/jpeg"}, nil
}

func (f *fakeSigner) Delete(_ context.Context, objectPath string) error {
	delete(f.objects, objectPath)
	return nil
}

// seedRestaurant creates a restaurant so the images FK is satisfied.
func seedRestaurant(ctx context.Context, db *dbtest.Database) uuid.UUID {
	orgs, err := organizationbus.TestSeedOrganizations(ctx, 1, db.BusDomain.Organization)
	if err != nil {
		panic(err)
	}
	rests, err := restaurantbus.TestSeedRestaurants(ctx, 1, db.BusDomain.Restaurant, orgs[0].ID)
	if err != nil {
		panic(err)
	}
	return rests[0].ID
}

func newTestBusiness(db *dbtest.Database, signer storage.Signer, maxSize int64) *imagebus.Business {
	return imagebus.NewBusiness(db.Log, imagedb.NewStore(db.Log, db.DB), signer, maxSize)
}

func Test_CreateUploadValidation(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_CreateUploadValidation")
	ctx := context.Background()
	restID := seedRestaurant(ctx, db)
	bus := newTestBusiness(db, newFakeSigner(), 1000)

	cases := []struct {
		name string
		req  imagebus.UploadRequest
	}{
		{"nil-restaurant", imagebus.UploadRequest{RestaurantID: uuid.Nil, EntityType: imagebus.EntityTypeRestaurant, ContentType: "image/jpeg", SizeBytes: 10}},
		{"bad-entity-type", imagebus.UploadRequest{RestaurantID: restID, EntityType: "addon", ContentType: "image/jpeg", SizeBytes: 10}},
		{"bad-content-type", imagebus.UploadRequest{RestaurantID: restID, EntityType: imagebus.EntityTypeMenuItem, ContentType: "image/gif", SizeBytes: 10}},
		{"zero-size", imagebus.UploadRequest{RestaurantID: restID, EntityType: imagebus.EntityTypeMenuItem, ContentType: "image/png", SizeBytes: 0}},
		{"negative-size", imagebus.UploadRequest{RestaurantID: restID, EntityType: imagebus.EntityTypeMenuItem, ContentType: "image/png", SizeBytes: -5}},
		{"oversize", imagebus.UploadRequest{RestaurantID: restID, EntityType: imagebus.EntityTypeMenuItem, ContentType: "image/png", SizeBytes: 1001}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := bus.CreateUpload(ctx, tc.req)
			if !errors.Is(err, imagebus.ErrInvalid) {
				t.Fatalf("got error %v, want ErrInvalid", err)
			}
		})
	}

	// Boundary: exactly the limit is allowed.
	grant, err := bus.CreateUpload(ctx, imagebus.UploadRequest{
		RestaurantID: restID,
		EntityType:   imagebus.EntityTypeMenuItem,
		ContentType:  "image/webp",
		SizeBytes:    1000,
	})
	if err != nil {
		t.Fatalf("boundary size should be accepted: %v", err)
	}
	if grant.Image.Status != imagebus.StatusPending {
		t.Errorf("status = %q, want pending", grant.Image.Status)
	}
	if !strings.HasPrefix(grant.Image.ObjectPath, "menu-items/"+restID.String()+"/") {
		t.Errorf("object path %q not scoped under menu-items/%s", grant.Image.ObjectPath, restID)
	}
	if !strings.HasSuffix(grant.Image.ObjectPath, ".webp") {
		t.Errorf("object path %q should end with .webp", grant.Image.ObjectPath)
	}
	if grant.UploadURL == "" || grant.ExpiresAt.IsZero() {
		t.Errorf("grant must carry upload url and expiry: %+v", grant)
	}
}

func Test_ConfirmUpload(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_ConfirmUpload")
	ctx := context.Background()
	restID := seedRestaurant(ctx, db)
	signer := newFakeSigner()
	bus := newTestBusiness(db, signer, 1000)

	grant, err := bus.CreateUpload(ctx, imagebus.UploadRequest{
		RestaurantID: restID,
		EntityType:   imagebus.EntityTypeRestaurant,
		ContentType:  "image/jpeg",
		SizeBytes:    500,
	})
	if err != nil {
		t.Fatalf("create upload: %v", err)
	}

	// Confirming before the object lands fails and leaves the row pending.
	if _, err := bus.ConfirmUpload(ctx, grant.Image.ID); !errors.Is(err, imagebus.ErrInvalid) {
		t.Fatalf("confirm before upload: got %v, want ErrInvalid", err)
	}
	img, err := bus.QueryByID(ctx, grant.Image.ID)
	if err != nil {
		t.Fatalf("query by id: %v", err)
	}
	if img.Status != imagebus.StatusPending {
		t.Errorf("status = %q, want pending after failed confirm", img.Status)
	}

	// Object lands; confirm flips status and adopts the actual size.
	signer.objects[grant.Image.ObjectPath] = 480

	img, err = bus.ConfirmUpload(ctx, grant.Image.ID)
	if err != nil {
		t.Fatalf("confirm: %v", err)
	}
	if img.Status != imagebus.StatusConfirmed {
		t.Errorf("status = %q, want confirmed", img.Status)
	}
	if img.SizeBytes != 480 {
		t.Errorf("size = %d, want actual 480", img.SizeBytes)
	}

	// Re-confirm is an idempotent no-op. Compare against the persisted row so
	// both sides carry database timestamp precision.
	persisted, err := bus.QueryByID(ctx, grant.Image.ID)
	if err != nil {
		t.Fatalf("query after confirm: %v", err)
	}
	img2, err := bus.ConfirmUpload(ctx, grant.Image.ID)
	if err != nil {
		t.Fatalf("re-confirm: %v", err)
	}
	if img2.Status != imagebus.StatusConfirmed || !img2.DateUpdated.Equal(persisted.DateUpdated) {
		t.Errorf("re-confirm should be a no-op, got %+v want date_updated %v", img2, persisted.DateUpdated)
	}

	// Unknown id surfaces ErrNotFound.
	if _, err := bus.ConfirmUpload(ctx, uuid.New()); !errors.Is(err, imagebus.ErrNotFound) {
		t.Fatalf("confirm unknown: got %v, want ErrNotFound", err)
	}
}

func Test_ConfirmUploadOversizeObject(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_ConfirmUploadOversizeObject")
	ctx := context.Background()
	restID := seedRestaurant(ctx, db)
	signer := newFakeSigner()
	bus := newTestBusiness(db, signer, 1000)

	grant, err := bus.CreateUpload(ctx, imagebus.UploadRequest{
		RestaurantID: restID,
		EntityType:   imagebus.EntityTypeMenuItem,
		ContentType:  "image/png",
		SizeBytes:    900,
	})
	if err != nil {
		t.Fatalf("create upload: %v", err)
	}

	// The client lied about the size: actual object exceeds the limit.
	signer.objects[grant.Image.ObjectPath] = 5000

	if _, err := bus.ConfirmUpload(ctx, grant.Image.ID); !errors.Is(err, imagebus.ErrInvalid) {
		t.Fatalf("confirm oversize: got %v, want ErrInvalid", err)
	}

	// Row and object must both be gone.
	if _, err := bus.QueryByID(ctx, grant.Image.ID); !errors.Is(err, imagebus.ErrNotFound) {
		t.Errorf("query after oversize confirm: got %v, want ErrNotFound", err)
	}
	if _, ok := signer.objects[grant.Image.ObjectPath]; ok {
		t.Error("oversize object should have been deleted from storage")
	}
}

func Test_Delete(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_Delete")
	ctx := context.Background()
	restID := seedRestaurant(ctx, db)
	signer := newFakeSigner()
	bus := newTestBusiness(db, signer, 0) // default limit

	grant, err := bus.CreateUpload(ctx, imagebus.UploadRequest{
		RestaurantID: restID,
		EntityType:   imagebus.EntityTypeRestaurant,
		ContentType:  "image/jpeg",
		SizeBytes:    100,
	})
	if err != nil {
		t.Fatalf("create upload: %v", err)
	}
	signer.objects[grant.Image.ObjectPath] = 100

	if err := bus.Delete(ctx, grant.Image.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := bus.QueryByID(ctx, grant.Image.ID); !errors.Is(err, imagebus.ErrNotFound) {
		t.Errorf("query after delete: got %v, want ErrNotFound", err)
	}
	if _, ok := signer.objects[grant.Image.ObjectPath]; ok {
		t.Error("object should have been deleted from storage")
	}

	if err := bus.Delete(ctx, uuid.New()); !errors.Is(err, imagebus.ErrNotFound) {
		t.Errorf("delete unknown: got %v, want ErrNotFound", err)
	}
}

func Test_Query(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_Query")
	ctx := context.Background()
	restID := seedRestaurant(ctx, db)
	otherRestID := seedRestaurant(ctx, db)
	signer := newFakeSigner()
	bus := newTestBusiness(db, signer, 0)

	upload := func(entityType string) imagebus.Image {
		grant, err := bus.CreateUpload(ctx, imagebus.UploadRequest{
			RestaurantID: restID,
			EntityType:   entityType,
			ContentType:  "image/jpeg",
			SizeBytes:    100,
		})
		if err != nil {
			t.Fatalf("create upload: %v", err)
		}
		return grant.Image
	}

	confirmed := upload(imagebus.EntityTypeRestaurant)
	signer.objects[confirmed.ObjectPath] = 100
	if _, err := bus.ConfirmUpload(ctx, confirmed.ID); err != nil {
		t.Fatalf("confirm: %v", err)
	}
	upload(imagebus.EntityTypeMenuItem) // stays pending

	// Unrelated restaurant upload must not leak into queries.
	otherGrant, err := bus.CreateUpload(ctx, imagebus.UploadRequest{
		RestaurantID: otherRestID,
		EntityType:   imagebus.EntityTypeRestaurant,
		ContentType:  "image/png",
		SizeBytes:    50,
	})
	if err != nil {
		t.Fatalf("create other upload: %v", err)
	}
	_ = otherGrant

	pg := page.MustParse("1", "10")

	all, err := bus.Query(ctx, imagebus.QueryFilter{RestaurantID: &restID}, imagebus.DefaultOrderBy, pg)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("query restaurant: got %d images, want 2", len(all))
	}

	status := imagebus.StatusConfirmed
	confirmedOnly, err := bus.Query(ctx, imagebus.QueryFilter{RestaurantID: &restID, Status: &status}, imagebus.DefaultOrderBy, pg)
	if err != nil {
		t.Fatalf("query confirmed: %v", err)
	}
	if len(confirmedOnly) != 1 || confirmedOnly[0].ID != confirmed.ID {
		t.Fatalf("query confirmed: got %+v, want only %s", confirmedOnly, confirmed.ID)
	}

	total, err := bus.Count(ctx, imagebus.QueryFilter{RestaurantID: &restID})
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	if total != 2 {
		t.Errorf("count = %d, want 2", total)
	}
}

// Test_DBTestDomainWiring ensures the shared test harness exposes the image domain.
func Test_DBTestDomainWiring(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_DBTestDomainWiring")
	if db.BusDomain.Image == nil {
		t.Fatal("BusDomain.Image should be wired")
	}
}
