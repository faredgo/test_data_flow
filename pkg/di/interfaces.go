package di

import (
	"archive/zip"
	assetschema "test_data_flow/internal/asset/schema"
	authschema "test_data_flow/internal/auth/schema"
	sessionschema "test_data_flow/internal/session/schema"
	userschema "test_data_flow/internal/user/schema"
)

// Repository interfaces
type IUserRepository interface {
	FindByLogin(login string) (*userschema.UserModel, error)
}

type IAssetRepository interface {
	Create(a *assetschema.AssetCommand) error
	Get(uid int64, assetName string) (*assetschema.AssetModel, error)
	Delete(uid int64, assetName string) error
	GetAll(uid int64) ([]*assetschema.AssetModel, error)
}

type ISessionRepository interface {
	DeleteByUID(uid int64) error
	Create(uid int64, ipAddress string) (string, error)
	GetByUID(uid int64) (*sessionschema.SessionModel, error)
}

// Service interfaces
type IAuthService interface {
	Login(loginCommand *authschema.LoginCommand) (int64, string, error)
}

type IAssetService interface {
	Upload(a *assetschema.AssetCommand) error
	Load(uid int64, assetName string) (*assetschema.AssetResponse, error)
	DeleteAsset(uid int64, assetName string) error
	GetAll(uid int64) ([]*assetschema.AssetResponse, error)
	MakeZip(zipWriter *zip.Writer, files []*assetschema.AssetResponse) error
}

type ISessionService interface {
	Delete(uid int64) error
	Create(uid int64, ipAddress string) (string, error)
	Get(uid int64) (*sessionschema.SessionResponse, error)
}
