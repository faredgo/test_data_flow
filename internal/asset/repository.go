package asset

import (
	"log"
	assetschema "test_data_flow/internal/asset/schema"

	"github.com/jmoiron/sqlx"
)

type AssetRepository struct {
	DB *sqlx.DB
}

func NewAssetRepository(db *sqlx.DB) *AssetRepository {
	return &AssetRepository{
		DB: db,
	}
}

func (r *AssetRepository) Create(assetCommand *assetschema.AssetCommand) error {
	model := &assetschema.AssetModel{
		Name: assetCommand.Name,
		UID:  assetCommand.UID,
		Data: assetCommand.Data,
	}

	query := `INSERT INTO assets (name, uid, data) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(query, model.Name, model.UID, model.Data)
	if err != nil {
		return err
	}

	return nil
}

func (r *AssetRepository) Get(uid int64, assetName string) (*assetschema.AssetModel, error) {
	var assetModel assetschema.AssetModel
	query := `SELECT name, uid, data, created_at FROM assets WHERE uid = $1 AND name = $2`
	err := r.DB.Get(&assetModel, query, uid, assetName)
	if err != nil {
		log.Printf("[REPO] Failed to get asset: %v", err)
		return nil, err
	}
	return &assetModel, nil
}

func (r *AssetRepository) Delete(uid int64, assetName string) error {
	query := `DELETE FROM assets WHERE uid = $1 AND name = $2`
	result, err := r.DB.Exec(query, uid, assetName)
	if err != nil {
		log.Printf("[REPO] Failed to delete asset: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[REPO] Failed to check affected rows: %v", err)
		return err
	}
	if rowsAffected == 0 {
		log.Println("No rows affected, asset not found")
		return nil
	}

	return nil
}

func (r *AssetRepository) GetAll(uid int64) ([]*assetschema.AssetModel, error) {
	var assets []*assetschema.AssetModel
	query := `SELECT name, uid, data, created_at FROM assets WHERE uid = $1`
	err := r.DB.Select(&assets, query, uid)
	if err != nil {
		log.Printf("[REPO] Failed to get all assets for uid %d: %v", uid, err)
		return nil, err
	}
	return assets, nil
}
