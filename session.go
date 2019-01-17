package influxdb

import (
	"context"
	"time"
)

// ErrSessionNotFound is the error messages for a missing sessions.
const ErrSessionNotFound = "session not found"

// ErrSessionExpired is the error message for expired sessions.
const ErrSessionExpired = "session has expired"

// RenewSessionTime is the the time to extend session, currently set to 5min.
var RenewSessionTime = time.Duration(time.Second * 300)

var (
	// OpFindSession represents the operation that looks for sessions.
	OpFindSession = "FindSession"
	// OpExpireSession represents the operation that expires sessions.
	OpExpireSession = "ExpireSession"
	// OpCreateSession represents the operation that creates a session for a given user.
	OpCreateSession = "CreateSession"
	// OpRenewSession = "RenewSession"
	OpRenewSession = "RenewSession"
)

// Session is a user session.
type Session struct {
	// ID is only required for auditing purposes.
	ID          ID           `json:"id"`
	Key         string       `json:"key"`
	CreatedAt   time.Time    `json:"createdAt"`
	ExpiresAt   time.Time    `json:"expiresAt"`
	UserID      ID           `json:"userID,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
}

// Expired returns an error if the session is expired.
func (s *Session) Expired() error {
	if time.Now().After(s.ExpiresAt) {
		return &Error{
			Code: EForbidden,
			Msg:  ErrSessionExpired,
		}
	}

	return nil
}

// Allowed returns true if the authorization is unexpired and request permission
// exists in the sessions list of permissions.
func (s *Session) Allowed(p Permission) bool {
	if err := s.Expired(); err != nil {
		return false
	}

	return PermissionAllowed(p, s.Permissions)
}

// Kind returns session and is used for auditing.
func (s *Session) Kind() string { return "session" }

// Identifier returns the sessions ID and is used for auditing.
func (s *Session) Identifier() ID { return s.ID }

// GetUserID returns the user id.
func (s *Session) GetUserID() ID {
	return s.UserID
}

// SessionService represents a service for managing user sessions.
type SessionService interface {
	FindSession(ctx context.Context, key string) (*Session, error)
	ExpireSession(ctx context.Context, key string) error
	CreateSession(ctx context.Context, user string) (*Session, error)
	RenewSession(ctx context.Context, session *Session, newExpiration time.Time) error
}
