package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/olzzhas/narxozer/graph/model"
	"github.com/olzzhas/narxozer/internal/validator"
	"golang.org/x/crypto/bcrypt"
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

func (u User) IsAnonymous() bool {
	return true
}

type password struct {
	plaintext *string
	hash      []byte
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type UserModel struct {
	DB    *sql.DB
	Redis *redis.Client
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

func (m UserModel) Insert(user *model.User) error {
	query := `
		INSERT INTO users (email, name, lastname, password_hash, role, image_url, additional_information, course, major, degree, faculty)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at`

	args := []interface{}{
		user.Email,
		user.Name,
		user.Lastname,
		user.PasswordHash,
		user.Role,
		user.ImageURL,
		user.AdditionalInformation,
		user.Course,
		user.Major,
		user.Degree,
		user.Faculty,
	}

	err := m.DB.QueryRow(query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (m UserModel) GetAll() ([]*model.User, error) {
	query := `
		SELECT id, email, name, lastname, role, image_url, additional_information, course, major, degree, faculty, created_at, updated_at
		FROM users
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Lastname,
			&user.Role,
			&user.ImageURL,
			&user.AdditionalInformation,
			&user.Course,
			&user.Major,
			&user.Degree,
			&user.Faculty,
			&user.CreatedAt,
			&user.UpdatedAt,
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

func (m UserModel) GetCached(userID int) (*model.User, error) {
	cacheKey := fmt.Sprintf("user:%d", userID)

	// Пытаемся получить данные из кеша Redis
	val, err := m.Redis.Get(context.Background(), cacheKey).Result()
	if err == redis.Nil {
		// Данные не найдены в кеше, загружаем из базы данных
		user, err := m.Get(userID)
		if err != nil {
			return nil, err
		}

		// Сохраняем данные в кеш Redis
		data, err := json.Marshal(user)
		if err != nil {
			return nil, err
		}
		err = m.Redis.Set(context.Background(), cacheKey, data, 10*time.Minute).Err()
		if err != nil {
			return nil, err
		}

		return user, nil
	} else if err != nil {
		return nil, err
	}

	// Если данные найдены в кеше, десериализуем их и возвращаем
	var user model.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m UserModel) Get(id int) (*model.User, error) {
	query := `
		SELECT id, email, name, lastname, role, image_url, additional_information, course, major, degree, faculty, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user model.User
	err := m.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Lastname,
		&user.Role,
		&user.ImageURL,
		&user.AdditionalInformation,
		&user.Course,
		&user.Major,
		&user.Degree,
		&user.Faculty,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (m UserModel) GetByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, email, name, lastname, password_hash, role, image_url, additional_information, course, major, degree, faculty, created_at, updated_at
		FROM users
		WHERE email = $1`

	var user model.User

	err := m.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Lastname,
		&user.PasswordHash,
		&user.Role,
		&user.ImageURL,
		&user.AdditionalInformation,
		&user.Course,
		&user.Major,
		&user.Degree,
		&user.Faculty,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (m UserModel) Update(id int, input model.UpdateUserInput) (*model.User, error) {
	query := `
		UPDATE users
		SET 
			email = COALESCE(NULLIF($1, ''), email),
			name = COALESCE(NULLIF($2, ''), name),
			lastname = COALESCE(NULLIF($3, ''), lastname),
			password_hash = COALESCE(NULLIF($4, ''), password_hash),
			role = COALESCE($5, role),
			image_url = COALESCE($6, image_url),
			additional_information = COALESCE($7, additional_information),
			course = COALESCE($8, course),
			major = COALESCE($9, major),
			degree = COALESCE($10, degree),
			faculty = COALESCE($11, faculty),
			updated_at = now()
		WHERE id = $12
		RETURNING id, email, name, lastname, password_hash, role, image_url, additional_information, course, major, degree, faculty, created_at, updated_at
	`

	user := &model.User{}
	err := m.DB.QueryRow(query,
		input.Email,
		input.Name,
		input.Lastname,
		input.Password,
		input.Role,
		input.ImageURL,
		input.AdditionalInformation,
		input.Course,
		input.Major,
		input.Degree,
		input.Faculty,
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Lastname,
		&user.PasswordHash,
		&user.Role,
		&user.ImageURL,
		&user.AdditionalInformation,
		&user.Course,
		&user.Major,
		&user.Degree,
		&user.Faculty,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user, nil
}
