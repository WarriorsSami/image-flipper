package main

import (
	"context"
	"image-flipper-tui/internal"
)

func main() {
	imageFolderPath := "/home/sami/Pictures/"
	ctx := context.Background()

	images, errChan := internal.ReadAllImagesInFolder(ctx, imageFolderPath)
}
