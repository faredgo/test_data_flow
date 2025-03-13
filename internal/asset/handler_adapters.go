package asset

import (
	assetschema "test_data_flow/internal/asset/schema"
)

func AssetFromRequest(name string, data []byte, uid int64) *assetschema.AssetCommand {
	return &assetschema.AssetCommand{
		Name: name,
		UID:  uid,
		Data: data,
	}
}
