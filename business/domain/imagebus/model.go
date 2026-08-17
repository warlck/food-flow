package imagebus

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Domain errors.
var (
	ErrNotFound = errors.New("image not found")
	ErrInvalid  = errors.New("invalid image data")
)

// Entity types an image can be associated with.
const (
	EntityTypeRestaurant = "restaurant"
	EntityTypeMenuItem   = "menu_item"
)

// Upload lifecycle statuses. Rows start as pending when the upload URL is
// minted and become confirmed once the object is verified in storage.
const (
	StatusPending   = "pending"
	StatusConfirmed = "confirmed"
)

// Image represents a tracked uploaded image object.
type Image struct {
	ID           uuid.UUID
	RestaurantID uuid.UUID
	EntityType   string
	ObjectPath   string
	PublicURL    string
	ContentType  string
	SizeBytes    int64
	Status       string
	UploadedBy   *uuid.UUID
	DateCreated  time.Time
	DateUpdated  time.Time
}
