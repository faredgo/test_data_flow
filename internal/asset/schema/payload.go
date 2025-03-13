package assetschema

type AssetResponse struct {
	Name string `json:"name"`
	File []byte `json:"file"`
}

type AssetCommand struct {
	Name string
	UID  int64
	Data []byte
}
