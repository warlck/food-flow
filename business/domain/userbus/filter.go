package userbus

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/name"
)

type QueryFilter struct {
	ID               *uuid.UUID
	Name             *name.Name
	Email            *mail.Address
	StartCreatedDate *time.Time
	EndCreatedDate   *time.Time
}
