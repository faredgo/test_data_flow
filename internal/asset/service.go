package asset

import (
	"archive/zip"
	"errors"
	"log"
	"mime"
	"net/http"
	assetschema "test_data_flow/internal/asset/schema"
	"test_data_flow/pkg/di"
)

type AssetService struct {
	AssetRepository di.IAssetRepository
}

func NewAssetService(assetRepository di.IAssetRepository) *AssetService {
	return &AssetService{
		AssetRepository: assetRepository,
	}
}

func (service *AssetService) Upload(asset *assetschema.AssetCommand) error {
	err := service.AssetRepository.Create(asset)
	if err != nil {
		log.Printf("[SRVC] Failed to create asset: %s", err)
		return errors.New("Failed to create asset")
	}

	return nil
}

func (service *AssetService) Load(uid int64, assetName string) (*assetschema.AssetResponse, error) {
	assetModel, err := service.AssetRepository.Get(uid, assetName)
	if err != nil {
		log.Printf("[SRVC] Failed to get asset: %s", err)
		return nil, errors.New("Failed to get asset")
	}

	return &assetschema.AssetResponse{
		Name: assetModel.Name,
		File: assetModel.Data,
	}, nil
}

func (service *AssetService) DeleteAsset(uid int64, assetName string) error {
	err := service.AssetRepository.Delete(uid, assetName)
	if err != nil {
		log.Printf("[SRVC] Failed to delete asset: %s", err)
		return errors.New("Failed to delete asset")
	}

	return nil
}

func (service *AssetService) GetAll(uid int64) ([]*assetschema.AssetResponse, error) {
	assetModels, err := service.AssetRepository.GetAll(uid)
	if err != nil {
		log.Printf("[SRVC] Failed to get all assets: %s", err)
		return nil, errors.New("Failed to get all assets")
	}

	var assetResponses []*assetschema.AssetResponse
	for _, assetModel := range assetModels {
		assetResponses = append(assetResponses, &assetschema.AssetResponse{
			Name: assetModel.Name,
			File: assetModel.Data,
		})
	}

	return assetResponses, nil
}

func (service *AssetService) MakeZip(zipWriter *zip.Writer, files []*assetschema.AssetResponse) error {
	for _, file := range files {
		fileType := http.DetectContentType(file.File)
		extensions, err := mime.ExtensionsByType(fileType)
		if err != nil || len(extensions) == 0 {
			log.Printf("[SRVC] Could not determine extension for MIME type: %s", fileType)
			extensions = []string{".bin"}
		}

		zipFileWriter, err := zipWriter.Create(file.Name + extensions[0])
		if err != nil {
			log.Printf("[SRVC] Error adding file to zip: %s", err)
			return errors.New("Error adding file to zip")
		}

		_, err = zipFileWriter.Write(file.File)
		if err != nil {
			log.Printf("[SRVC] Error writing file to zip: %s", err)
			return errors.New("Error writing file to zip")
		}
	}

	return nil
}
