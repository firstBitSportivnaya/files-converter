package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/firstBitSportivnaya/files-converter/pkg/converter"
	"github.com/spf13/cobra"
)

var (
	// srcConvertCmd represents the srcConvert command
	srcConvertCmd = &cobra.Command{
		Use:   "srcConvert",
		Short: "Tool for converting source files to *.cfe.",
		Run: func(cmd *cobra.Command, args []string) {
			runConverter()
		},
	}
	inputPath  string
	outputPath string
)

func init() {
	rootCmd.AddCommand(srcConvertCmd)

	srcConvertCmd.Flags().StringVarP(&inputPath, "input", "i", "", "Source directory")
	srcConvertCmd.MarkFlagRequired("input")
	srcConvertCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output directory")
	srcConvertCmd.MarkFlagRequired("output")
}

func runConverter() {
	defer pressAnyKeyToExit()

	sourceDir := filepath.Clean(strings.TrimSpace(inputPath))
	targetDir := filepath.Clean(strings.TrimSpace(outputPath))

	if err := converter.ConvertToCfe(sourceDir, targetDir); err != nil {
		log.Printf("Could not to convert files: %v", err)
	}
}

func pressAnyKeyToExit() {
	fmt.Println("Press any key to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
