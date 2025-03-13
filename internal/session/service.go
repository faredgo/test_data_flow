package session

import (
	"log"
	sessionschema "test_data_flow/internal/session/schema"
	"test_data_flow/pkg/di"
)

type SessionService struct {
	SessionRepository di.ISessionRepository
}

func NewSessionService(sessionRepository di.ISessionRepository) *SessionService {
	return &SessionService{
		SessionRepository: sessionRepository,
	}
}

func (s *SessionService) Delete(uid int64) error {
	err := s.SessionRepository.DeleteByUID(uid)
	if err != nil {
		log.Printf("[SRVC] Failed to delete session for user ID %d: %s", uid, err)
		return err
	}
	return nil
}

func (s *SessionService) Create(uid int64, ipAddress string) (string, error) {
	sessionID, err := s.SessionRepository.Create(uid, ipAddress)
	if err != nil {
		log.Printf("[SRVC] Failed to create session for user ID %d from IP %s: %s", uid, ipAddress, err)
		return "", err
	}
	return sessionID, nil
}

func (s *SessionService) Get(uid int64) (*sessionschema.SessionResponse, error) {
	sessionResp, err := s.SessionRepository.GetByUID(uid)
	if err != nil {
		log.Printf("[SRVC] Failed to get session for user ID %d: %s", uid, err)
		return nil, err
	}
	return &sessionschema.SessionResponse{
		ID:        sessionResp.ID,
		UID:       sessionResp.UID,
		CreatedAt: sessionResp.CreatedAt,
	}, nil
}
