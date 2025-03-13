package user

import (
	"log"
	userschema "test_data_flow/internal/user/schema"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (repo *UserRepository) FindByLogin(login string) (*userschema.UserModel, error) {
	var user userschema.UserModel

	query := `SELECT id, login, password_hash, created_at FROM users WHERE login = $1 LIMIT 1`

	err := repo.DB.Get(&user, query, login)
	if err != nil {
		log.Printf("[REPO] Failed to find user with login %s: %s", login, err)
		return nil, err
	}

	return &user, nil
}
