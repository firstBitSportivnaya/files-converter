package cmd

import (
	"log"

	"github.com/firstBitSportivnaya/files-converter/pkg/config"
	"github.com/firstBitSportivnaya/files-converter/pkg/converter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	// временно
	viper.SetConfigFile("configs/config.json")
	cfg, err := config.LoadConfig(viper.GetViper())
	if err != nil {
		log.Println("Ошибка загрузки конфигурации:", err)
		return
	}
	cfg.InputPath = sourcePath
	cfg.OutputPath = targetDir
	// временно
	if err := converter.ConvertFromCf(cfg, sourcePath, targetDir); err != nil {
		log.Printf("Could not to convert files: %v", err)
	}
}
