package cmd

import (
	"log"

	"github.com/firstBitSportivnaya/files-converter/pkg/config"
	"github.com/firstBitSportivnaya/files-converter/pkg/converter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	// временно
	viper.SetConfigFile("configs/config.json")
	cfg, err := config.LoadConfig(viper.GetViper())
	if err != nil {
		log.Println("Ошибка загрузки конфигурации:", err)
		return
	}
	cfg.InputPath = sourceDir
	cfg.OutputPath = targetDir
	// временно
	if err := converter.ConvertFromSourceFiles(cfg, sourceDir, targetDir); err != nil {
		log.Printf("Could not to convert files: %v", err)
	}
}
