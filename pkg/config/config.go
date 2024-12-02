package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	_ "embed"
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

func LoadConfig() (*Configuration, error) {
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Configuration
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("не удалось разобрать конфигурацию: %w", err)
	}

	if err := normalizeConfigPaths(&config); err != nil {
		return nil, err
	}

	return &config, nil
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

func GetDefaultConfig() (*Configuration, error) {
	var config Configuration

	err := json.Unmarshal(defaultConfig, &config)
	if err != nil {
		return nil, fmt.Errorf("не удалось разобрать конфигурацию по умолчанию: %w", err)
	}

	if err := normalizeConfigPaths(&config); err != nil {
		return nil, err
	}

	return &config, nil
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
