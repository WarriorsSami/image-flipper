package internal

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

func runReadAllImagesFromFolderStage(ctx context.Context, g *errgroup.Group, folderPath string) <-chan *Image {
	readImgChan := make(chan *Image)
	var wg sync.WaitGroup

	g.Go(func() error {
		defer close(readImgChan)

		files, err := filepath.Glob(folderPath + "/*[.png|.jpg|.jpeg|.bmp]")
		if err != nil {
			return err
		}

		errChan := make(chan error)
		defer close(errChan)
		for _, file := range files {
			if !isImageFile(file) {
				continue
			}

			wg.Add(1)
			go func(file string) {
				defer wg.Done()
				select {
				case <-ctx.Done():
					return
				default:
					img := NewImage(file)
					img, err := readImage(img)
					if err != nil {
						errChan <- err
						return
					}

					log.Println("Read image:", img.name)
					readImgChan <- img
				}
			}(file)
		}
		wg.Wait()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errChan:
			return err
		default:
			return nil
		}
	})

	return readImgChan
}

func runFlipImagesStage(ctx context.Context, g *errgroup.Group, images <-chan *Image, direction FlipDirection) <-chan *Image {
	flippedImgChan := make(chan *Image)
	var wg sync.WaitGroup

	g.Go(func() error {
		defer close(flippedImgChan)

		errChan := make(chan error)
		defer close(errChan)
		for img := range images {
			wg.Add(1)
			go func(img *Image) {
				defer wg.Done()

				select {
				case <-ctx.Done():
					return
				default:
					flippedImg, err := flipImage(img, direction)
					if err != nil {
						errChan <- err
						return
					}

					log.Println("Flipped image:", img.name)
					flippedImgChan <- flippedImg
				}
			}(img)
		}
		wg.Wait()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errChan:
			return err
		default:
			return nil
		}
	})

	return flippedImgChan
}

func runWriteImagesToFolderStage(ctx context.Context, g *errgroup.Group, images <-chan *Image, outputFolderPath string) <-chan *Image {
	writtenImgChan := make(chan *Image)
	var wg sync.WaitGroup

	g.Go(func() error {
		defer close(writtenImgChan)

		errChan := make(chan error)
		defer close(errChan)
		for img := range images {
			wg.Add(1)
			go func(img *Image) {
				defer wg.Done()
				select {
				case <-ctx.Done():
					return
				default:
					if err := writeImage(img, outputFolderPath); err != nil {
						errChan <- err
						return
					}

					log.Println("Wrote image:", img.name)
					writtenImgChan <- img
				}
			}(img)
		}
		wg.Wait()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errChan:
			return err
		default:
			return nil
		}
	})

	return writtenImgChan
}

func RunProcessImagesPipeline(ctx context.Context, inputImagesFolderPath, outputImagesFolderPath string, direction FlipDirection) (string, error) {
	g, ctx := errgroup.WithContext(ctx)

	readImagesChan := runReadAllImagesFromFolderStage(ctx, g, inputImagesFolderPath)
	flippedImagesChan := runFlipImagesStage(ctx, g, readImagesChan, direction)
	writtenImagesChan := runWriteImagesToFolderStage(ctx, g, flippedImagesChan, outputImagesFolderPath)

	var successMsgs strings.Builder

	g.Go(func() error {
		for img := range writtenImagesChan {
			successMsgs.WriteString(fmt.Sprintf("Successfully processed image %s\n", img.name))
		}
		successMsgs.WriteString("All images processed successfully")

		return nil
	})

	if err := g.Wait(); err != nil {
		log.Println("Error processing images:", err)
		return "", err
	}

	return successMsgs.String(), nil
}
