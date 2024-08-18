// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Comment struct {
	ID        string     `json:"id"`
	Content   string     `json:"content"`
	AuthorID  string     `json:"authorId"`
	CreatedAt string     `json:"createdAt"`
	Likes     int        `json:"likes"`
	Replies   []*Comment `json:"replies"`
}

type CreateCommentInput struct {
	PostID   string `json:"postId"`
	Content  string `json:"content"`
	AuthorID string `json:"authorId"`
}

type CreatePostInput struct {
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	ImageURL *string `json:"imageURL,omitempty"`
	AuthorID string  `json:"authorId"`
}

type CreateUserInput struct {
	Email                 string  `json:"email"`
	Name                  string  `json:"name"`
	Lastname              string  `json:"lastname"`
	Password              string  `json:"password"`
	Role                  Role    `json:"role"`
	ImageURL              *string `json:"imageURL,omitempty"`
	AdditionalInformation *string `json:"additionalInformation,omitempty"`
	Course                *int    `json:"course,omitempty"`
	Major                 *string `json:"major,omitempty"`
	Degree                *string `json:"degree,omitempty"`
	Faculty               *string `json:"faculty,omitempty"`
}

type Mutation struct {
}

type Post struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	ImageURL  *string    `json:"imageURL,omitempty"`
	AuthorID  string     `json:"authorId"`
	CreatedAt string     `json:"createdAt"`
	UpdatedAt *string    `json:"updatedAt,omitempty"`
	Likes     int        `json:"likes"`
	Comments  []*Comment `json:"comments"`
}

type Query struct {
}

type UpdatePostInput struct {
	Title    *string `json:"title,omitempty"`
	Content  *string `json:"content,omitempty"`
	ImageURL *string `json:"imageURL,omitempty"`
}

type UpdateUserInput struct {
	Email                 *string `json:"email,omitempty"`
	Name                  *string `json:"name,omitempty"`
	Lastname              *string `json:"lastname,omitempty"`
	Password              *string `json:"password,omitempty"`
	Role                  *Role   `json:"role,omitempty"`
	ImageURL              *string `json:"imageURL,omitempty"`
	AdditionalInformation *string `json:"additionalInformation,omitempty"`
	Course                *int    `json:"course,omitempty"`
	Major                 *string `json:"major,omitempty"`
	Degree                *string `json:"degree,omitempty"`
	Faculty               *string `json:"faculty,omitempty"`
}

type User struct {
	ID                    string  `json:"id"`
	Email                 string  `json:"email"`
	Name                  string  `json:"name"`
	Lastname              string  `json:"lastname"`
	PasswordHash          string  `json:"passwordHash"`
	Role                  Role    `json:"role"`
	ImageURL              *string `json:"imageURL,omitempty"`
	AdditionalInformation *string `json:"additionalInformation,omitempty"`
	Course                *int    `json:"course,omitempty"`
	Major                 *string `json:"major,omitempty"`
	Degree                *string `json:"degree,omitempty"`
	Faculty               *string `json:"faculty,omitempty"`
}

type Role string

const (
	RoleStudent Role = "STUDENT"
	RoleTeacher Role = "TEACHER"
	RoleAdmin   Role = "ADMIN"
)

var AllRole = []Role{
	RoleStudent,
	RoleTeacher,
	RoleAdmin,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleStudent, RoleTeacher, RoleAdmin:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
