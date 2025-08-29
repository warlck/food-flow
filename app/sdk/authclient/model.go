package authclient

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/auth"
)

// Authorize defines the information required to perform an authorization.
type Authorize struct {
	UserID uuid.UUID
	Claims auth.Claims
	Rule   string
}

// Decode implements the decoder interface.
func (a *Authorize) Decode(data []byte) error {
	return json.Unmarshal(data, a)
}

// AuthenticateResp defines the information that will be received on authenticate.
type AuthenticateResp struct {
	UserID uuid.UUID
	Claims auth.Claims
}
