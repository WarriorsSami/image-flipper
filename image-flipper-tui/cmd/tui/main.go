package main

import (
	"context"
	imgproc "image_utils"
	"log"
)

func main() {
	imageFolderPath, outputFolderPath := "/home/sami/Pictures/", "/home/sami/Pictures/output/"
	flipDir := imgproc.FlipHorizontal
	ctx := context.Background()

	successMessage, err := imgproc.RunProcessImagesPipeline(ctx, imageFolderPath, outputFolderPath, flipDir)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(successMessage)
}
