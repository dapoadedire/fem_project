package store

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plaintext *string
	hash      []byte
}

type User struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	PasswordHash   password  `json:"password_hash"`
	Bio            string    `json:"bio"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	ProfilePicture string    `json:"profile_picture"`
	LastLogin      string    `json:"last_login"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
	Workouts       []Workout `json:"workouts"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	CreateUser(*User) (*User, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUser(*User) error
	DeleteUser(id int64) error
	GetUserByEmail(email string) (*User, error)
}

func (s *PostgresUserStore) CreateUser(user *User) (*User, error) {
	query := `INSERT INTO users(username, email, password_hash, bio, first_name, last_name, profile_picture)
	VALUES($1,$2,$3,$4,$5,$6,$7)
	RETURNING id, created_at, updated_at
	`
	err := s.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.hash, user.Bio, user.FirstName, user.LastName, user.ProfilePicture).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}
