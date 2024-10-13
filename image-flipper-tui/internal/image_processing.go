package internal

import (
	"errors"
	"fmt"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"image"
	"path/filepath"
)

// FlipDirection - enum for image flip direction
type FlipDirection int

const (
	FlipHorizontal FlipDirection = iota
	FlipVertical
	FlipBoth
)

// ImgExtension - enum for image extension
type ImgExtension int

const (
	PNG ImgExtension = iota
	JPEG
	BMP
)

func (f *FlipDirection) String() string {
	switch *f {
	case FlipHorizontal:
		return "horizontal"
	case FlipVertical:
		return "vertical"
	case FlipBoth:
		return "both"
	default:
		return "unknown"
	}
}

func (f *FlipDirection) Set(value string) error {
	switch value {
	case "horizontal":
		*f = FlipHorizontal
	case "vertical":
		*f = FlipVertical
	case "both":
		*f = FlipBoth
	default:
		return errors.New("invalid flip direction")
	}

	return nil
}

func (f *FlipDirection) Type() string {
	return "string"
}

type Image struct {
	path      string
	name      string
	extension ImgExtension
	image     *image.Image
}

func NewImage(path string) *Image {
	name := filepath.Base(path)
	var extension ImgExtension
	switch filepath.Ext(path) {
	case ".png":
		extension = PNG
	case ".jpg", ".jpeg":
		extension = JPEG
	case ".bmp":
		extension = BMP
	default:
		extension = PNG
	}

	return &Image{
		path:      path,
		name:      name,
		extension: extension,
		image:     nil,
	}
}

func (img *Image) ToShortString() string {
	return fmt.Sprintf("Image{name: %s}", img.name)
}

func (img *Image) getImageEncoder() imgio.Encoder {
	switch img.extension {
	case PNG:
		return imgio.PNGEncoder()
	case JPEG:
		return imgio.JPEGEncoder(100)
	case BMP:
		return imgio.BMPEncoder()
	default:
		return imgio.PNGEncoder()
	}
}

func isImageFile(path string) bool {
	switch filepath.Ext(path) {
	case ".png", ".jpg", ".jpeg", ".bmp":
		return true
	default:
		return false
	}
}

func readImage(img *Image) (*Image, error) {
	imgRaw, err := imgio.Open(img.path)
	if err != nil {
		return nil, err
	}
	img.image = &imgRaw

	return img, nil
}

func flipImage(img *Image, direction FlipDirection) (*Image, error) {
	var flippedImage *image.RGBA
	switch direction {
	case FlipHorizontal:
		flippedImage = transform.FlipH(*img.image)
	case FlipVertical:
		flippedImage = transform.FlipV(*img.image)
	case FlipBoth:
		flippedImage = transform.FlipH(*img.image)
		flippedImage = transform.FlipV(flippedImage)
	}

	if flippedImage == nil {
		return nil, errors.New("error flipping image")
	}

	newImg := flippedImage.SubImage(flippedImage.Bounds())
	img.image = &newImg

	return img, nil
}

func writeImage(img *Image, outputFolderPath string) error {
	outputPath := filepath.Join(outputFolderPath, img.name)
	err := imgio.Save(outputPath, *img.image, img.getImageEncoder())
	if err != nil {
		return err
	}

	return nil
}

func CheckIfFolderExists(folderPath string) (bool, error) {
	if _, err := filepath.Abs(folderPath); err != nil {
		return false, err
	}

	return true, nil
}
