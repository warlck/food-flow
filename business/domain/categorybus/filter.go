package categorybus

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/name"
)

// QueryFilter holds the available fields for filtering Category queries.
type QueryFilter struct {
	ID               *uuid.UUID
	Name             *name.Name
	RestaurantID     *uuid.UUID
	Enabled          *bool
	StartCreatedDate *time.Time
	EndCreatedDate   *time.Time
}
