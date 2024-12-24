package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "embed"

	"github.com/ilyakaznacheev/cleanenv"
)

//go:embed default.json
var defaultConfig []byte

var globalDumpInfo *DumpInfo

func init() {
	globalDumpInfo = NewDumpInfo()
}

func NewDumpInfo() *DumpInfo {
	return &DumpInfo{ConfigName: "default"}
}

func GetDumpInfo() *DumpInfo {
	return globalDumpInfo
}

func (info *DumpInfo) SetVersion(version string) {
	if version != "" {
		info.Version = "V" + version
	}
}

func (info *DumpInfo) SetConfigName(name string) {
	if name != "" {
		info.ConfigName = name
	}
}

func LoadConfig(configPath string) *Configuration {
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	if configPath == "" {
		panic("не указан путь конфигурации")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("конфигурационный файл не существует: " + err.Error())
	}

	var cfg Configuration

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("ошибка чтения конфигурационного файла: " + err.Error())
	}

	if err := normalizeConfigPaths(&cfg); err != nil {
		panic("ошибка в обработке путей: " + err.Error())
	}

	return &cfg
}

func LoadDefaultConfig() *Configuration {
	var cfg Configuration

	if err := json.Unmarshal(defaultConfig, &cfg); err != nil {
		panic("не удалось прочитать конфигурацию по умолчанию: " + err.Error())
	}

	if err := normalizeConfigPaths(&cfg); err != nil {
		panic("ошибка в обработке путей: " + err.Error())
	}

	return &cfg
}

func NormalizePath(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	expandedPath := os.ExpandEnv(strings.TrimSpace(input))

	absPath, err := filepath.Abs(expandedPath)
	if err != nil {
		return "", err
	}

	cleanPath := filepath.Clean(absPath)

	return cleanPath, nil
}

func NewElementOperation(name string, value string, operation OperationType) *ElementOperation {
	return &ElementOperation{
		ElementName: name,
		Value:       value,
		Operation:   operation,
	}
}

func normalizeConfigPaths(config *Configuration) error {
	var err error
	config.InputPath, err = NormalizePath(config.InputPath)
	if err != nil {
		return fmt.Errorf("не удалось нормализовать путь к входному файлу: %w", err)
	}

	config.OutputPath, err = NormalizePath(config.OutputPath)
	if err != nil {
		return fmt.Errorf("не удалось нормализовать путь к выходному файлу: %w", err)
	}

	return nil
}
