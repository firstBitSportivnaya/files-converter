package cmd

import (
	"log"

	"github.com/firstBitSportivnaya/files-converter/pkg/converter"
	"github.com/spf13/cobra"
)

// srcConvertCmd represents the srcConvert command
var srcConvertCmd = &cobra.Command{
	Use:   "srcConvert",
	Short: "Tool for converting source files to *.cfe.",
	Run: func(cmd *cobra.Command, args []string) {
		runConverterSrc()
	},
}

func init() {
	rootCmd.AddCommand(srcConvertCmd)

	srcConvertCmd.Flags().StringVarP(&inputPath, "input", "i", "", "Source directory")
	srcConvertCmd.MarkFlagRequired("input")
	srcConvertCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output directory")
	srcConvertCmd.MarkFlagRequired("output")
}

func runConverterSrc() {
	defer pressAnyKeyToExit()

	sourceDir, targetDir := NormalizePaths(inputPath, outputPath)

	if err := converter.ConvertFromSourceFiles(sourceDir, targetDir); err != nil {
		log.Printf("Could not to convert files: %v", err)
	}
}
