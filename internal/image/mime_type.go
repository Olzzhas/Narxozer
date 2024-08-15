package image

var validImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

func IsAllowedImageType(mimeType string) bool {
	_, exist := validImageTypes[mimeType]

	return exist
}
