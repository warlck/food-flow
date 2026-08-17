package imagedb

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/imagebus"
)

type image struct {
	ID           uuid.UUID  `db:"image_id"`
	RestaurantID uuid.UUID  `db:"restaurant_id"`
	EntityType   string     `db:"entity_type"`
	ObjectPath   string     `db:"object_path"`
	PublicURL    string     `db:"public_url"`
	ContentType  string     `db:"content_type"`
	SizeBytes    int64      `db:"size_bytes"`
	Status       string     `db:"status"`
	UploadedBy   *uuid.UUID `db:"uploaded_by"`
	DateCreated  time.Time  `db:"date_created"`
	DateUpdated  time.Time  `db:"date_updated"`
}

func toDBImage(bus imagebus.Image) image {
	return image{
		ID:           bus.ID,
		RestaurantID: bus.RestaurantID,
		EntityType:   bus.EntityType,
		ObjectPath:   bus.ObjectPath,
		PublicURL:    bus.PublicURL,
		ContentType:  bus.ContentType,
		SizeBytes:    bus.SizeBytes,
		Status:       bus.Status,
		UploadedBy:   bus.UploadedBy,
		DateCreated:  bus.DateCreated.UTC(),
		DateUpdated:  bus.DateUpdated.UTC(),
	}
}

func toBusImage(db image) imagebus.Image {
	return imagebus.Image{
		ID:           db.ID,
		RestaurantID: db.RestaurantID,
		EntityType:   db.EntityType,
		ObjectPath:   db.ObjectPath,
		PublicURL:    db.PublicURL,
		ContentType:  db.ContentType,
		SizeBytes:    db.SizeBytes,
		Status:       db.Status,
		UploadedBy:   db.UploadedBy,
		DateCreated:  db.DateCreated.In(time.Local),
		DateUpdated:  db.DateUpdated.In(time.Local),
	}
}

func toBusImages(dbs []image) []imagebus.Image {
	bus := make([]imagebus.Image, len(dbs))
	for i, db := range dbs {
		bus[i] = toBusImage(db)
	}
	return bus
}
