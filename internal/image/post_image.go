package image

import "cloud.google.com/go/storage"

type PostImageStorage struct {
	Client *storage.Client
}
