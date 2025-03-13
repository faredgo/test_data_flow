package asset_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"test_data_flow/internal/asset"
	assetschema "test_data_flow/internal/asset/schema"
	"test_data_flow/pkg/middleware"
	mockdi "test_data_flow/testmocks/pkg/di"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	UID      int64  = 1
	FileName string = "dummyx"
)

func TestAssetHandlers(t *testing.T) {
	mockAssetRepository := mockdi.NewMockIAssetRepository(t)
	assetService := asset.NewAssetService(mockAssetRepository)
	handler := asset.AssetHandler{
		AssetService: assetService,
	}

	cwd, err := os.Getwd()
	require.NoError(t, err)

	filePath := filepath.Join(cwd, "../../dummyx.pdf")
	dataFile, err := os.ReadFile(filePath)
	require.NoError(t, err)
	t.Run("Upload Asset Success", func(t *testing.T) {
		mockAssetRepository.On("Create", mock.Anything).Return(nil).Once()

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		part, err := writer.CreateFormFile("file", filepath.Base(FileName))
		require.NoError(t, err)
		_, err = part.Write(dataFile)
		require.NoError(t, err)
		require.NoError(t, writer.Close())

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/upload-asset/%s", FileName), &body)
		r.Header.Set("Content-Type", writer.FormDataContentType())

		ctx := context.WithValue(r.Context(), middleware.ContextKeyID, UID)
		r = r.WithContext(ctx)
		handler.Upload()(w, r)

		require.Equal(t, http.StatusOK, w.Code)

		mockAssetRepository.AssertExpectations(t)
	})

	t.Run("Load Asset Success", func(t *testing.T) {
		mockAssetRepository.On("Get", UID, mock.Anything).Return(&assetschema.AssetModel{
			Name:      FileName,
			UID:       UID,
			Data:      dataFile,
			CreatedAt: time.Now(),
		}, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/asset/%s", FileName), nil)

		ctx := context.WithValue(r.Context(), middleware.ContextKeyID, UID)
		r = r.WithContext(ctx)
		handler.Load()(w, r)

		require.Equal(t, http.StatusOK, w.Code)

		mockAssetRepository.AssertExpectations(t)
	})

	t.Run("Delete Asset Success", func(t *testing.T) {
		mockAssetRepository.On("Delete", UID, mock.Anything).Return(nil).Maybe()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/asset/%s", FileName), nil)

		ctx := context.WithValue(r.Context(), middleware.ContextKeyID, UID)
		r = r.WithContext(ctx)
		handler.Load()(w, r)

		require.Equal(t, http.StatusOK, w.Code)

		mockAssetRepository.AssertExpectations(t)
	})

	t.Run("Load All Assets Success", func(t *testing.T) {
		mockAssetRepository.On("GetAll", UID).Return([]*assetschema.AssetModel{
			{
				UID:       UID,
				Name:      "asset1",
				Data:      dataFile,
				CreatedAt: time.Now(),
			},
			{
				UID:       UID,
				Name:      "asset2",
				Data:      dataFile,
				CreatedAt: time.Now(),
			},
		}, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/assets", nil)

		ctx := context.WithValue(r.Context(), middleware.ContextKeyID, UID)
		r = r.WithContext(ctx)
		handler.LoadAll()(w, r)

		require.Equal(t, http.StatusOK, w.Code)

		mockAssetRepository.AssertExpectations(t)
	})

	t.Run("Upload Asset Failure - Invalid File", func(t *testing.T) {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		_, err := writer.CreateFormFile("file", "")
		require.NoError(t, err)
		require.NoError(t, writer.Close())

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/upload-asset/%s", FileName), &body)
		r.Header.Set("Content-Type", writer.FormDataContentType())

		ctx := context.WithValue(r.Context(), middleware.ContextKeyID, UID)
		r = r.WithContext(ctx)
		handler.Upload()(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Upload Asset Failure - Repository Error", func(t *testing.T) {
		mockAssetRepository.On("Create", mock.Anything).Return(errors.New("failed to save file"))

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		part, err := writer.CreateFormFile("file", filepath.Base(FileName))
		require.NoError(t, err)
		_, err = part.Write(dataFile)
		require.NoError(t, err)
		require.NoError(t, writer.Close())

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/upload-asset/%s", FileName), &body)
		r.Header.Set("Content-Type", writer.FormDataContentType())

		ctx := context.WithValue(r.Context(), middleware.ContextKeyID, UID)
		r = r.WithContext(ctx)
		handler.Upload()(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code)

		mockAssetRepository.AssertExpectations(t)
	})
}
