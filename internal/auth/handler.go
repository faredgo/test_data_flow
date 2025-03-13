package auth

import (
	"encoding/json"
	"net/http"
	"test_data_flow/configs"
	authschema "test_data_flow/internal/auth/schema"
	"test_data_flow/pkg/di"
	"test_data_flow/pkg/jwt"
	"test_data_flow/pkg/req"
	"test_data_flow/pkg/res"
)

const (
	postAuth string = "POST /api/auth"
)

type AuthHandlerDeps struct {
	Config      *configs.Config
	AuthService di.IAuthService
}

type AuthHandler struct {
	Config      *configs.Config
	AuthService di.IAuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc(postAuth, handler.Login())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload authschema.LoginRequest
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}

		// TODO: validation
		if payload.Login == "" {
			res.ReturnError(w, "Login required", http.StatusBadRequest)
			return
		}
		if payload.Password == "" {
			res.ReturnError(w, "Password required", http.StatusBadRequest)
			return
		}

		ipAddress := req.GetIPAddress(r)
		loginCommand := &authschema.LoginCommand{
			Login:     payload.Login,
			Password:  payload.Password,
			IpAdderss: ipAddress,
		}
		userID, sessionID, err := handler.AuthService.Login(loginCommand)
		if err != nil {
			res.ReturnError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(&jwt.JWTData{
			ID:        userID,
			Login:     payload.Login,
			SessionID: sessionID,
		})
		if err != nil {
			res.ReturnError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, authschema.LoginResponse{
			Token: token,
		}, http.StatusOK)
	}
}
