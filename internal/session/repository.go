package session

import (
	"log"
	sessionschema "test_data_flow/internal/session/schema"

	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	DB *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{
		DB: db,
	}
}

func (r *SessionRepository) DeleteByUID(uid int64) error {
	_, err := r.DB.Exec(`DELETE FROM sessions WHERE uid = $1`, uid)
	if err != nil {
		log.Printf("[REPO] Failed to delete session for user ID %d: %s", uid, err)
		return err
	}
	return nil
}

func (r *SessionRepository) Create(uid int64, ipAddress string) (string, error) {
	var sessionID string
	query := `INSERT INTO sessions (uid, ip_address) VALUES ($1, $2) RETURNING id`
	err := r.DB.QueryRow(query, uid, ipAddress).Scan(&sessionID)
	if err != nil {
		log.Printf("[REPO] Failed to create session for user ID %d with IP %s: %s", uid, ipAddress, err)
		return "", err
	}

	return sessionID, nil
}

func (r *SessionRepository) GetByUID(uid int64) (*sessionschema.SessionModel, error) {
	var session sessionschema.SessionModel
	err := r.DB.Get(&session, `SELECT id, uid, ip_address, created_at FROM sessions WHERE uid = $1`, uid)
	if err != nil {
		log.Printf("[REPO] Failed to retrieve session for user ID %d: %s", uid, err)
		return nil, err
	}
	return &session, nil
}
