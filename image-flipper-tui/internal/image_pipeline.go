package internal

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

func readAllImagesFromFolder(ctx context.Context, done chan struct{}, folderPath string) (<-chan *Image, chan error) {
	readImgChan := make(chan *Image)
	errChan := make(chan error)
	var wg sync.WaitGroup

	go func() {
		defer close(readImgChan)

		// TODO: refactor using filepath.Walk
		files, err := filepath.Glob(folderPath + "/*")
		if err != nil {
			errChan <- err
			done <- struct{}{}
			return
		}

		for _, file := range files {
			if !isImageFile(file) {
				continue
			}

			wg.Add(1)
			go func(file string) {
				defer wg.Done()
				select {
				case <-ctx.Done():
				case <-done:
					return
				default:
					img := NewImage(file)
					img, err := readImage(img)
					if err != nil {
						errChan <- err
						return
					}

					readImgChan <- img
				}
			}(file)
		}

		wg.Wait()
	}()

	return readImgChan, errChan
}

func flipImages(ctx context.Context, done <-chan struct{}, images <-chan *Image, errChan chan error, direction FlipDirection) (<-chan *Image, chan error) {
	flippedImgChan := make(chan *Image)
	var wg sync.WaitGroup

	go func() {
		defer close(flippedImgChan)

		for img := range images {
			wg.Add(1)
			go func(img *Image) {
				defer wg.Done()

				select {
				case <-ctx.Done():
				case <-done:
					return
				default:
					flippedImg, err := flipImage(img, direction)
					if err != nil {
						errChan <- err
						return
					}

					flippedImgChan <- flippedImg
				}
			}(img)
		}

		wg.Wait()
	}()

	return flippedImgChan, errChan
}

func writeImagesToFolder(ctx context.Context, done chan struct{}, images <-chan *Image, errChan chan error, outputFolderPath string) (<-chan *Image, <-chan error) {
	writtenImgChan := make(chan *Image)
	var wg sync.WaitGroup

	go func() {
		defer close(writtenImgChan)
		defer close(errChan)
		defer func() {
			done <- struct{}{}
		}()

		for img := range images {
			wg.Add(1)
			go func(img *Image) {
				defer wg.Done()
				select {
				case <-ctx.Done():
				case <-done:
					return
				default:
					if err := writeImage(img, outputFolderPath); err != nil {
						errChan <- err
						return
					}

					writtenImgChan <- img
				}
			}(img)
		}

		wg.Wait()
	}()

	return writtenImgChan, errChan
}

// TODO: refactor using errgroup
func RunProcessImagesPipeline(ctx context.Context, inputImagesFolderPath, outputImagesFolderPath string, direction FlipDirection) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	done := make(chan struct{})

	readImagesChan, readErrChan := readAllImagesFromFolder(ctx, done, inputImagesFolderPath)
	flippedImagesChan, flipErrChan := flipImages(ctx, done, readImagesChan, readErrChan, FlipVertical)
	writtenImagesChan, writeErrChan := writeImagesToFolder(ctx, done, flippedImagesChan, flipErrChan, outputImagesFolderPath)

	var (
		successMsgs strings.Builder
		errMsgs     strings.Builder
	)

	for {
		select {
		case img := <-writtenImagesChan:
			if img != nil {
				log.Println("Wrote image:", img.ToShortString())
			}
			successMsgs.WriteString(fmt.Sprintf("Wrote image: %s\n", img.ToShortString()))
		case err := <-writeErrChan:
			if err != nil {
				log.Println("Error writing image:", err)
			}
			errMsgs.WriteString(fmt.Sprintf("Error writing image: %s\n", err))
		case <-done:
			log.Println("Done processing images")
			close(done)

			successMsgs.WriteString("Done processing images\n")
			if errMsgs.Len() == 0 {
				return successMsgs.String(), nil
			}
			return successMsgs.String(), fmt.Errorf(errMsgs.String())
		case <-ctx.Done():
			log.Println("Processing images cancelled")
			close(done)

			successMsgs.WriteString("Processing images cancelled\n")
			if errMsgs.Len() == 0 {
				return successMsgs.String(), nil
			}
			return successMsgs.String(), fmt.Errorf(errMsgs.String())
		}
	}
}
