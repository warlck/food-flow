package restaurantbus

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/name"
)

// QueryFilter holds the available fields for filtering Restaurant queries.
type QueryFilter struct {
	ID               *uuid.UUID
	OrganizationID   *uuid.UUID
	Name             *name.Name
	Enabled          *bool
	StartCreatedDate *time.Time
	EndCreatedDate   *time.Time
}
