package auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"test_data_flow/configs"
	"test_data_flow/internal/auth"
	authschema "test_data_flow/internal/auth/schema"
	"test_data_flow/internal/session"
	userschema "test_data_flow/internal/user/schema"
	mockdi "test_data_flow/testmocks/pkg/di"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	UID       int64  = 1
	IpAddress string = "192.0.2.1"
	SessionID string = "5ebe2294ecd0e0f08eab7690d2a6ee70"
)

func TestLogin(t *testing.T) {
	mockUserRepository := mockdi.NewMockIUserRepository(t)
	mockSessionRepository := mockdi.NewMockISessionRepository(t)
	sessionService := session.NewSessionService(mockSessionRepository)
	authService := auth.NewAuthService(mockUserRepository, sessionService)
	mockSessionRepository.On("DeleteByUID", UID).Return(nil)
	mockSessionRepository.On("Create", UID, IpAddress).Return(SessionID, nil)
	handler := auth.AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "secret",
			},
		},
		AuthService: authService,
	}

	t.Run("Success", func(t *testing.T) {
		loginRequest := &authschema.LoginRequest{
			Login:    "alice",
			Password: "secret",
		}

		mockUserRepository.
			On("FindByLogin", loginRequest.Login).
			Return(&userschema.UserModel{
				ID:        1,
				Login:     "alice",
				Password:  "5ebe2294ecd0e0f08eab7690d2a6ee69",
				CreatedAt: time.Now(),
			}, nil)

		data, _ := json.Marshal(loginRequest)
		reader := bytes.NewReader(data)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/auth", reader)

		handler.Login()(w, r)

		require.Equal(t, http.StatusOK, w.Code)

		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Invalid password", func(t *testing.T) {
		loginRequest := &authschema.LoginRequest{
			Login:    "alice",
			Password: "wrong-password",
		}

		mockUserRepository.ExpectedCalls = nil
		mockUserRepository.
			On("FindByLogin", loginRequest.Login).
			Return(&userschema.UserModel{
				ID:        1,
				Login:     "alice",
				Password:  "5ebe2294ecd0e0f08eab7690d2a6ee69",
				CreatedAt: time.Now(),
			}, nil)

		data, _ := json.Marshal(loginRequest)
		reader := bytes.NewReader(data)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/auth", reader)

		handler.Login()(w, r)

		require.Equal(t, http.StatusUnauthorized, w.Code)

		mockUserRepository.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		loginRequest := &authschema.LoginRequest{
			Login:    "bob",
			Password: "111",
		}

		mockUserRepository.ExpectedCalls = nil
		mockUserRepository.
			On("FindByLogin", loginRequest.Login).
			Return(nil, errors.New("user not found"))

		data, _ := json.Marshal(loginRequest)
		reader := bytes.NewReader(data)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/auth", reader)

		handler.Login()(w, r)

		require.Equal(t, http.StatusUnauthorized, w.Code)

		mockUserRepository.AssertExpectations(t)
	})
}
