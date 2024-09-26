package utils

import "strings"

type (
	ImageTag struct {
		image string
		tag   string
	}
)

// ImageParser parses the image name and tag
func ImageParser(image string) ImageTag {
	x := ImageTag{
		image: image,
	}

	imageParsed := strings.Split(image, ":")
	if len(imageParsed) > 1 {
		x.image = imageParsed[0]
		x.tag = imageParsed[1]
	}

	return x
}

// GetImage returns the image name
func (i ImageTag) GetImage() string {
	return i.image
}

// TagExists checks if the tag exists
func (i ImageTag) TagExists() bool {
	return i.tag != ""
}

// GetImageWithoutTag returns the image name without the tag
func (i ImageTag) GetImageWithoutTag() string {
	return i.GetImage()
}

// GetImageWithTag returns the image name with the tag
func (i ImageTag) GetImageWithTag() string {
	return i.image + ":" + i.tag
}

// GetTag returns the tag
func (i ImageTag) GetTag() string {
	return i.tag
}
