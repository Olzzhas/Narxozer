package data

import (
	"cloud.google.com/go/storage"
	"github.com/olzzhas/narxozer/internal/image"
)

type Storages struct {
	ProfileImage image.ProfileImageStorage
	PostImage    image.PostImageStorage
}

func NewStorages(client *storage.Client) Storages {
	return Storages{
		ProfileImage: image.ProfileImageStorage{Client: client},
		PostImage:    image.PostImageStorage{Client: client},
	}
}
