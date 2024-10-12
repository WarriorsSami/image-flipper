package main

import (
	"context"
	imgproc "image_utils"
	"log"
)

func main() {
	imageFolderPath, outputFolderPath := "/home/sami/Pictures/", "/home/sami/Pictures/output/"
	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, "cancel", cancel)

	images, errChan := imgproc.ReadAllImagesInFolder(ctx, imageFolderPath)
	flippedImages, flipErrChan := imgproc.FlipImages(ctx, images, errChan, imgproc.FlipBoth)
	writeErrChan := imgproc.WriteImagesToFolder(ctx, flippedImages, flipErrChan, outputFolderPath)

	for {
		select {
		case err := <-writeErrChan:
			if err != nil {
				log.Println("Error writing image:", err)
			}
		case <-ctx.Done():
			log.Println("Done processing images")
			return
		}
	}
}
