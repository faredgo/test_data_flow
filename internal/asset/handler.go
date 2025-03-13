package asset

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"test_data_flow/configs"
	"test_data_flow/pkg/di"
	"test_data_flow/pkg/middleware"
	"test_data_flow/pkg/res"
)

const (
	postAsset   string = "POST /api/upload-asset/{name}"
	getAsset    string = "GET /api/asset/{name}"
	deleteAsset string = "DELETE /api/asset/{name}"
	getAllAsset string = "GET /api/assets"
)

type AssetHandlerDeps struct {
	AssetService di.IAssetService
	Config       *configs.Config
}

type AssetHandler struct {
	AssetService di.IAssetService
}

func NewAssetHandler(router *http.ServeMux, deps AssetHandlerDeps) {
	handler := &AssetHandler{
		AssetService: deps.AssetService,
	}
	router.Handle(postAsset, middleware.IsAuthed(handler.Upload(), deps.Config))
	router.Handle(getAsset, middleware.IsAuthed(handler.Load(), deps.Config))
	router.Handle(deleteAsset, middleware.IsAuthed(handler.Delete(), deps.Config))
	router.Handle(getAllAsset, middleware.IsAuthed(handler.LoadAll(), deps.Config))
}

func (handler *AssetHandler) Upload() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id, ok := req.Context().Value(middleware.ContextKeyID).(int64)
		if !ok {
			res.ReturnError(w, "Login not found", http.StatusUnauthorized)
			return
		}

		file, _, err := req.FormFile("file") // header.Filename
		if err != nil {
			log.Printf("Failed to get file: %s", err.Error())
			res.ReturnError(w, "Failed to get file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			log.Printf("Failed to read file: %s", err.Error())
			res.ReturnError(w, "Failed to read file", http.StatusInternalServerError)
			return
		}
		assetName := req.PathValue("name")
		assetCommand := AssetFromRequest(assetName, data, id)
		err = handler.AssetService.Upload(assetCommand)
		if err != nil {
			log.Printf("Failed to upload file: %s", err.Error())
			res.ReturnError(w, "Error uploading asset", http.StatusInternalServerError)
			return
		}

		res.Json(w, "ок", http.StatusOK)
	}
}

func (handler *AssetHandler) Load() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		assetName := req.PathValue("name")
		id, ok := req.Context().Value(middleware.ContextKeyID).(int64)
		if !ok {
			res.ReturnError(w, "Login not found", http.StatusUnauthorized)
			return
		}

		assetResp, err := handler.AssetService.Load(id, assetName)
		if err != nil {
			log.Printf("Failed to getting file: %s", err.Error())
			res.ReturnError(w, "Error getting asset", http.StatusInternalServerError)
			return
		}

		fileType := http.DetectContentType(assetResp.File)
		extensions, err := mime.ExtensionsByType(fileType)
		if err != nil || len(extensions) == 0 {
			log.Printf("Could not determine extension for MIME type: %s", fileType)
			res.ReturnError(w, "Error could not determine extension for MIME type", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+assetName+extensions[0])
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(assetResp.File)))

		w.Write(assetResp.File)
	}
}

func (handler *AssetHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id, ok := req.Context().Value(middleware.ContextKeyID).(int64)
		if !ok {
			res.ReturnError(w, "Login not found", http.StatusUnauthorized)
			return
		}

		assetName := req.PathValue("name")

		err := handler.AssetService.DeleteAsset(id, assetName)
		if err != nil {
			log.Printf("Failed to delete file: %s", err.Error())
			res.ReturnError(w, "Error delete asset", http.StatusInternalServerError)
			return
		}
	}
}

func (handler *AssetHandler) LoadAll() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id, ok := req.Context().Value(middleware.ContextKeyID).(int64)
		if !ok {
			res.ReturnError(w, "Login not found", http.StatusUnauthorized)
			return
		}

		files, err := handler.AssetService.GetAll(id)
		if err != nil {
			log.Printf("Failed to get all files: %s", err.Error())
			res.ReturnError(w, "Error getting assets", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=assets.zip")
		zipWriter := zip.NewWriter(w)
		defer zipWriter.Close()

		err = handler.AssetService.MakeZip(zipWriter, files)
		if err != nil {
			log.Printf("Failed to make zip: %s", err.Error())
			res.ReturnError(w, "Error making zip", http.StatusInternalServerError)
			return
		}
	}
}
