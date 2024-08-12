package cmd

import (
	"log"

	"github.com/firstBitSportivnaya/files-converter/pkg/converter"
	"github.com/spf13/cobra"
)

// cfConvertCmd represents the cfConvert command
var cfConvertCmd = &cobra.Command{
	Use:   "cfConvert",
	Short: "Tool for converting file *.cf to *.cfe.",
	Run: func(cmd *cobra.Command, args []string) {
		runConverterCf()
	},
}

func init() {
	rootCmd.AddCommand(cfConvertCmd)

	cfConvertCmd.Flags().StringVarP(&inputPath, "input", "i", "", "Source directory")
	cfConvertCmd.MarkFlagRequired("input")
	cfConvertCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output directory")
	cfConvertCmd.MarkFlagRequired("output")
}

func runConverterCf() {
	defer pressAnyKeyToExit()

	sourcePath, targetDir := NormalizePaths(inputPath, outputPath)

	if err := converter.ConvertFromCf(sourcePath, targetDir); err != nil {
		log.Printf("Could not to convert files: %v", err)
	}
}
