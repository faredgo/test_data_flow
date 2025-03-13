package assetschema

import "time"

type AssetModel struct {
	Name      string    `db:"name"`
	UID       int64     `db:"uid"`
	Data      []byte    `db:"data"`
	CreatedAt time.Time `db:"created_at"`
}
