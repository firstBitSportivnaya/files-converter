package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/firstBitSportivnaya/files-converter/pkg/config"
	"github.com/firstBitSportivnaya/files-converter/pkg/converter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "files-converter",
	Short: "A tool for converting files to the *.cfe format",
	Long: `Files Converter is a command-line application that allows you to convert files to *.cfe format.
	
There are two conversion modes available:
1. Convert from source files to *.cfe.
2. Convert from .cf file to *.cfe.

This tool simplifies the conversion process, making it easy and efficient to manage your files.`,
	Run: runMain,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.files-converter/configs/config.json)")

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("json")
		viper.AddConfigPath("./configs")

		viper.AddConfigPath("$HOME/.files-converter/configs")
		viper.AddConfigPath("/etc/files-converter")
	}
	viper.AutomaticEnv()
}

func runMain(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfig(viper.GetViper())
	if err != nil {
		log.Fatalf("ошибка загрузки конфигурации: %v", err)
	}

	fmt.Println("Используется файл конфигурации:", viper.ConfigFileUsed())

	runConvert(cfg)
}

func runConvert(cfg *config.Configuration) {
	defer pressAnyKeyToExit()

	if err := converter.RunConversion(cfg); err != nil {
		log.Fatalf("не удалось конвертировать файлы: %v", err)
	}
}

func pressAnyKeyToExit() {
	fmt.Println("Press any key to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
