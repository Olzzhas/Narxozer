package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/olzzhas/narxozer/internal/validator"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

var AnonymousUser = &User{}

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Role      string    `json:"role"`
	ImageUrl  string    `json:"image_url"`
	LFKUrl    string    `json:"lfk_url"`
	LFKAccess bool      `json:"lfk_access"`
	Version   int       `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type UserModel struct {
	DB *sql.DB
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(user.Surname != "", "surname", "must be provided")
	v.Check(len(user.Surname) <= 500, "surname", "must not be more than 500 bytes long")

	//v.Check(user.Role == "admin" || user.Role == "student" || user.Role == "teacher", "role", "invalid user role")

	ValidateEmail(v, user.Email)

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (name, surname, email, password_hash, activated, role, image_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, version, role
	`

	args := []any{user.Name, user.Surname, user.Email, user.Password.hash, user.Activated, "student", ""}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version, &user.Role)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, created_at, name, surname, email, password_hash, activated, version, role, image_url, lfk_url, lfk_access
		FROM users
		WHERE email = $1
	`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Surname,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
		&user.Role,
		&user.ImageUrl,
		&user.LFKUrl,
		&user.LFKAccess,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) Get(id int64) (*User, error) {
	query := `
		SELECT id, created_at, name, surname, email, password_hash, activated, version, role, image_url, lfk_url, lfk_access
		FROM users
		WHERE id = $1
	`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Surname,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
		&user.Role,
		&user.ImageUrl,
		&user.LFKUrl,
		&user.LFKAccess,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET name = $1, surname = $2, email = $3, password_hash = $4, activated = $5, version = version + 1, role = $8, image_url = $9
		WHERE id = $6 AND version = $7
		RETURNING version
	`

	args := []any{
		user.Name,
		user.Surname,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
		user.Role,
		user.ImageUrl,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (m UserModel) CheckUserExists(id int64) (bool, error) {
	sqlStmt := `SELECT 1 FROM users WHERE id = $1`

	stmt, err := m.DB.Prepare(sqlStmt)
	if err != nil {
		return false, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			return
		}
	}(stmt)

	row := stmt.QueryRow(id)

	var exists bool
	err = row.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (m UserModel) FindAllByLesson(id int64) ([]*User, error) {
	query := `
		SELECT users.id, users.created_at, users.name, users.surname, users.email, users.activated, users.role, users.image_url, users.version, users.lfk_url, users.lfk_access
        FROM lesson_registrations
        JOIN users ON users.id = lesson_registrations.student_id
        WHERE lesson_registrations.lesson_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	//defer rows.Close()
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var users []*User

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.CreatedAt,
			&user.Name,
			&user.Surname,
			&user.Email,
			&user.Activated,
			&user.Role,
			&user.ImageUrl,
			&user.Version,
			&user.LFKUrl,
			&user.LFKAccess,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
