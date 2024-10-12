package internal

import (
	"context"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"image"
	"path/filepath"
	"sync"
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

type Image struct {
	path      string
	name      string
	extension ImgExtension
	image     *image.RGBA
}

func NewImageMeta(path string) *Image {
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

func (img *Image) GetImageEncoder() imgio.Encoder {
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

func ReadAllImagesInFolder(ctx context.Context, folderPath string) (<-chan *Image, <-chan error) {
	readImgChan := make(chan *Image)
	errorChan := make(chan error)
	var wg sync.WaitGroup

	go func() {
		defer close(readImgChan)
		defer close(errorChan)

		files, err := filepath.Glob(folderPath + "/*")
		if err != nil {
			errorChan <- err
			ctx.Done()
			return
		}

		for _, file := range files {
			if !isImageFile(file) {
				continue
			}

			wg.Add(1)
			go func(file string) {
				defer wg.Done()
				img := NewImageMeta(file)
				img, err := readImage(img)
				if err != nil {
					errorChan <- err
					return
				}
				readImgChan <- img
			}(file)
		}

		wg.Wait()
	}()

	return readImgChan, errorChan
}

func FlipImages(ctx context.Context, images <-chan *Image, errChan chan error, direction FlipDirection) (<-chan *Image, <-chan error) {
	flippedImgChan := make(chan *Image)
	var wg sync.WaitGroup

	go func() {
		defer close(flippedImgChan)
		for img := range images {
			wg.Add(1)
			go func(img *Image) {
				defer wg.Done()
				flippedImg, err := flipImage(img, direction)
				if err != nil {
					errChan <- err
					return
				}
				flippedImgChan <- flippedImg
			}(img)
		}
	}()
	return flippedImgChan, errChan
}

func WriteImagesToFolder(ctx context.Context, images <-chan *Image, errChan chan error, folderPath string) <-chan error {
	errorChan := make(chan error)
	var wg sync.WaitGroup

	go func() {
		defer close(errorChan)
		for img := range images {
			wg.Add(1)
			go func(img *Image) {
				defer wg.Done()
				if err := writeImage(img); err != nil {
					errorChan <- err
					return
				}
			}(img)
		}
		wg.Wait()
	}()

	return errorChan
}

func readImage(img *Image) (*Image, error) {
	imgRaw, err := imgio.Open(img.path)
	if err != nil {
		return nil, err
	}
	img.image = imgRaw.(*image.RGBA)

	return img, nil
}

func flipImage(img *Image, direction FlipDirection) (*Image, error) {
	var flippedImage *image.RGBA
	switch direction {
	case FlipHorizontal:
		flippedImage = transform.FlipH(img.image)
	case FlipVertical:
		flippedImage = transform.FlipV(img.image)
	case FlipBoth:
		flippedImage = transform.FlipH(img.image)
		flippedImage = transform.FlipV(flippedImage)
	}

	img.image = flippedImage
	return img, nil
}

func writeImage(img *Image) error {
	err := imgio.Save(img.path, img.image, img.GetImageEncoder())
	if err != nil {
		return err
	}

	return nil
}
