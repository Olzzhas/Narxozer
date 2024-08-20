package data

import (
	"database/sql"
	"errors"
	"github.com/go-redis/redis/v8"
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

func NewModels(db *sql.DB, redis *redis.Client) Models {
	return Models{
		Permissions:         PermissionModel{DB: db, Redis: redis},
		Users:               UserModel{DB: db, Redis: redis},
		Tokens:              TokenModel{DB: db, Redis: redis},
		AuthorizationTokens: AuthorizationTokenModel{DB: db, Redis: redis},
		Posts:               PostModel{DB: db, Redis: redis},
		Clubs:               ClubModel{DB: db, Redis: redis},
		Events:              EventModel{DB: db, Redis: redis},
		Topics:              TopicModel{DB: db, Redis: redis},
		Comments:            CommentModel{DB: db, Redis: redis},
	}
}
