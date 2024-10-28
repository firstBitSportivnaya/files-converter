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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.files-converter/configs/config.json)")

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func runMain(cmd *cobra.Command, args []string) {
	initDefaultConfig()
	defaultCfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("ошибка загрузки конфигурации: %v", err)
	}

	initConfig()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("ошибка загрузки конфигурации: %v", err)
	}

	mergeConfigs(defaultCfg, cfg)

	changeXmlFiles(defaultCfg)

	runConvert(defaultCfg)
}

func changeXmlFiles(cfg *config.Configuration) {
	for _, xmlFile := range cfg.XMLFiles {
		if xmlFile.FileName == "Configuration.xml" {
			setNamePrefix(xmlFile, cfg.Prefix)
		}
	}
}

func setNamePrefix(file *config.FileOperation, prefix string) {
	for _, operation := range file.ElementOperations {
		if operation.ElementName == "NamePrefix" {
			operation.Value = prefix
			return
		}
	}
	element := &config.ElementOperation{
		ElementName: "NamePrefix",
		Value:       prefix,
		Operation:   config.Add,
	}
	file.ElementOperations = append(file.ElementOperations, element)
}

func initDefaultConfig() {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("default")
	viper.SetConfigType("json")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("json")
		viper.AddConfigPath("./configs")

		viper.AddConfigPath("$HOME/.files-converter/configs")
		viper.AddConfigPath("/etc/files-converter/configs")
	}
}

func mergeConfigs(defaultCfg, cfg *config.Configuration) {
	if cfg.PlatformVersion != "" {
		defaultCfg.PlatformVersion = cfg.PlatformVersion
	}

	defaultCfg.Extension = cfg.Extension
	defaultCfg.Prefix = cfg.Prefix
	defaultCfg.InputPath = cfg.InputPath
	defaultCfg.OutputPath = cfg.OutputPath
	defaultCfg.ConversionType = cfg.ConversionType

	defaultCfg.XMLFiles = append(defaultCfg.XMLFiles, cfg.XMLFiles...)
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
