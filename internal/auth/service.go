package auth

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	authschema "test_data_flow/internal/auth/schema"
	"test_data_flow/pkg/di"
)

type AuthService struct {
	UserRepository di.IUserRepository
	SessionService di.ISessionService
}

func NewAuthService(userRepository di.IUserRepository, sessionService di.ISessionService) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
		SessionService: sessionService,
	}
}

func (service *AuthService) Login(loginCommand *authschema.LoginCommand) (int64, string, error) {
	user, err := service.UserRepository.FindByLogin(loginCommand.Login)
	if err != nil {
		log.Printf("[SRVC] User get error for login %s: %s", loginCommand.Login, err)
		return -1, "", errors.New(ErrWrongCredentials)
	}

	md5Hash := md5.Sum([]byte(loginCommand.Password))
	md5String := hex.EncodeToString(md5Hash[:])
	if md5String != user.Password {
		log.Printf("[SRVC] Incorrect password for user %s", loginCommand.Login)
		return -1, "", errors.New(ErrWrongCredentials)
	}

	err = service.SessionService.Delete(user.ID)
	if err != nil {
		log.Printf("[SRVC] Failed to delete existing session for user %s: %s", loginCommand.Login, err)
		return -1, "", err
	}
	sessionID, err := service.SessionService.Create(user.ID, loginCommand.IpAdderss)
	if err != nil {
		log.Printf("[SRVC] Failed to create session for user %s: %s", loginCommand.Login, err)
		return -1, "", err
	}

	return user.ID, sessionID, nil
}
