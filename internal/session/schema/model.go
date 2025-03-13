package sessionschema

import (
	"time"
)

type SessionModel struct {
	ID        string    `db:"id"`
	UID       int64     `db:"uid"`
	IpAddress string    `db:"ip_address"`
	CreatedAt time.Time `db:"created_at"`
}
