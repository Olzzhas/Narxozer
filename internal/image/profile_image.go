package image

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
)

type ProfileImageStorage struct {
	Client *storage.Client
}

func (st ProfileImageStorage) UpdateProfile(ctx context.Context, objName string, imageFile multipart.File) (string, error) {
	profileImageBucketName := os.Getenv("GC_PROFILE_IMAGE_BUCKET")
	bucket := st.Client.Bucket(profileImageBucketName)

	object := bucket.Object(objName)
	wc := object.NewWriter(ctx)

	wc.ObjectAttrs.CacheControl = "Cache-Control:no-cache, max-age=0"

	if _, err := io.Copy(wc, imageFile); err != nil {
		log.Printf("Unable to write file to Google Cloud Storage: %v\n", err)
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	imageURL := fmt.Sprintf(
		"https://storage.googleapis.com/%s/%s",
		profileImageBucketName,
		objName,
	)

	return imageURL, nil
}
