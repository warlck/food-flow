package imagebus

import (
	"github.com/google/uuid"
)

// QueryFilter holds fields for filtering Image queries.
type QueryFilter struct {
	ID           *uuid.UUID
	RestaurantID *uuid.UUID
	EntityType   *string
	Status       *string
}
