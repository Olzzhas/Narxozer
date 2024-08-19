package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Permissions         PermissionModel
	Users               UserModel
	Tokens              TokenModel
	AuthorizationTokens AuthorizationTokenModel
	Posts               PostModel
	Clubs               ClubModel
	Events              EventModel
	Topics              TopicModel
	Comments            CommentModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Permissions:         PermissionModel{DB: db},
		Users:               UserModel{DB: db},
		Tokens:              TokenModel{DB: db},
		AuthorizationTokens: AuthorizationTokenModel{DB: db},
		Posts:               PostModel{DB: db},
		Clubs:               ClubModel{DB: db},
		Events:              EventModel{DB: db},
		Topics:              TopicModel{DB: db},
		Comments:            CommentModel{DB: db},
	}
}
