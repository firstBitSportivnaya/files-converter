package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var info *DumpInfo

func init() {
	info = New()
}

func (info *DumpInfo) SetVersion(in string) {
	if in != "" {
		info.Version = "V" + in
	}
}

func (info *DumpInfo) SetConfigName(in string) {
	if in != "" {
		info.ConfigName = in
	}
}

func New() *DumpInfo {
	info := new(DumpInfo)
	info.ConfigName = "default"

	return info
}

func GetDumpInfo() *DumpInfo {
	return info
}

func LoadConfig() (*Configuration, error) {
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Configuration
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	path, err := NormalizePath(config.InputPath)
	if err != nil {
		return nil, err
	}
	config.InputPath = path

	path, err = NormalizePath(config.OutputPath)
	if err != nil {
		return nil, err
	}
	config.OutputPath = path

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
