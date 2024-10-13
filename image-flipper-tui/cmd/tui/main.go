package main

import (
	"context"
	imgproc "image_utils"
	"log"
)

func main() {
	imageFolderPath, outputFolderPath := "/home/sami/Pictures/", "/home/sami/Pictures/output/"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})

	readImagesChan, readErrChan := imgproc.ReadAllImagesInFolder(ctx, done, imageFolderPath)
	flippedImagesChan, flipErrChan := imgproc.FlipImages(ctx, done, readImagesChan, readErrChan, imgproc.FlipVertical)
	writtenImagesChan, writeErrChan := imgproc.WriteImagesToFolder(ctx, done, flippedImagesChan, flipErrChan, outputFolderPath)

	for {
		select {
		case img := <-writtenImagesChan:
			if img != nil {
				log.Println("Wrote image:", img.ToShortString())
			}
		case err := <-writeErrChan:
			if err != nil {
				log.Println("Error writing image:", err)
			}
		case <-done:
			log.Println("Done processing images")
			close(done)
			return
		}
	}
}
