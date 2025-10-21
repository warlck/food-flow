package menuitembus

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
)

// QueryFilter holds the available fields for filtering MenuItem queries.
type QueryFilter struct {
	ID               *uuid.UUID
	Name             *name.Name
	CategoryID       *uuid.UUID
	RestaurantID     *uuid.UUID
	MinPrice         *money.Money
	MaxPrice         *money.Money
	Available        *bool
	StartCreatedDate *time.Time
	EndCreatedDate   *time.Time
}
