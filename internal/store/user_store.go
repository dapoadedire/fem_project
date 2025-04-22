package store

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

func (password *password) SetPassword(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	password.plainText = &plaintextPassword
	password.hash = hash
	return nil
}

func (password *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(password.hash, []byte(plaintextPassword))
	
	if err != nil {
		switch {
			case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
				return false, nil
			default:
				return false, err // internal server error
		}
	}
	return true, nil
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
	LastLogin      *string    `json:"last_login"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
	Workouts       []Workout `json:"workouts"`
}

var AnonymousUser = &User{} // EVERYONE WHOS NOT LOGGED IN
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	CreateUser(*User) ( error)
	GetUserByUsername(username string) (*User, error)
	UpdateUser(*User) error
	GetUserToken(scope, tokenPlainText string) (*User, error)
	// DeleteUser(id int64) error
	// GetUserByEmail(email string) (*User, error)
}

func (s *PostgresUserStore) CreateUser(user *User) error {
	query := `INSERT INTO users(username, email, password_hash, bio, first_name, last_name, profile_picture)
	VALUES($1,$2,$3,$4,$5,$6,$7)
	RETURNING id, created_at, updated_at
	`
	err := s.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.hash, user.Bio, user.FirstName, user.LastName, user.ProfilePicture).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}
	query := `SELECT id, username, email, password_hash, bio, first_name, last_name, profile_picture, last_login, created_at, updated_at
	FROM users WHERE username = $1`

	err := s.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash.hash,
		&user.Bio, &user.FirstName, &user.LastName, &user.ProfilePicture,
		&user.LastLogin, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *PostgresUserStore) UpdateUser(user *User) error {
	query := `UPDATE users SET username = $1, email = $2, password_hash = $3, bio = $4, first_name = $5, last_name = $6, profile_picture = $7, updated_at = CURRENT_TIMESTAMP
	WHERE id = $8
	RETURNING updated_at
	`
	result, err := s.db.Exec(query, user.Username, user.Email, user.PasswordHash.hash, user.Bio, user.FirstName, user.LastName, user.ProfilePicture, user.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}



func (s *PostgresUserStore) GetUserToken(scope, plaintextPassword string) (*User, error) {
	rokenHash := sha256.Sum256([]byte(plaintextPassword))
	
	query := `SELECT id, username, email, password_hash, bio, first_name, last_name, profile_picture, last_login, created_at, updated_at
	FROM users u
	INNER JOIN tokens t ON t.user_id = u.id
	WHERE t.hash = $1 AND t.scope = $2 and t.expiry > $3
	
	`

	user := &User{
		PasswordHash: password{},
	}
	err := s.db.QueryRow(query, rokenHash[:], scope, time.Now()).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash.hash,
		&user.Bio, &user.FirstName, &user.LastName, &user.ProfilePicture,
		&user.LastLogin, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}