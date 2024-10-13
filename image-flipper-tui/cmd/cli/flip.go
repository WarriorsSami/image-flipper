package main

import (
	"context"
	"errors"
	"github.com/spf13/cobra"
	imgproc "image_utils"
	"log"
)

var (
	direction = imgproc.FlipHorizontal
	flipCmd   = &cobra.Command{
		Use:   "flip",
		Short: "Flip images",
		Long:  "Flip images horizontally, vertically or both",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 0 {
				return errors.New("requires at least one argument")
			}

			input, err := cmd.Flags().GetString("input")
			if err != nil {
				return err
			}
			output, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}
			flipDirection, err := cmd.Flags().GetString("direction")
			if err != nil {
				return err
			}

			// check if input/output folders exist
			if _, err := imgproc.CheckIfFolderExists(input); err != nil {
				return err
			}
			if _, err := imgproc.CheckIfFolderExists(output); err != nil {
				return err
			}

			// check if flip direction is valid
			if err := direction.Set(flipDirection); err != nil {
				return err
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			input, _ := cmd.Flags().GetString("input")
			output, _ := cmd.Flags().GetString("output")

			log.Printf("Input: %s, Output: %s, Direction: %s\n", input, output, direction.String())
			ctx := context.Background()

			successMessage, err := imgproc.RunProcessImagesPipeline(ctx, input, output, direction)
			if err != nil {
				log.Fatalln(err)
			}

			log.Println(successMessage)
		},
	}
)

func Execute() {
	if err := flipCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	flipCmd.Flags().StringP("input", "i", "", "Input folder path")
	flipCmd.Flags().StringP("output", "o", "", "Output folder path")
	flipCmd.Flags().VarP(&direction, "direction", "d", "Flip direction. Possible values: horizontal, vertical, both")
}
