package sessionschema

import "time"

type SessionResponse struct {
	ID        string    `json:"id"`
	UID       int64     `json:"uid"`
	CreatedAt time.Time `json:"created_at"`
}
