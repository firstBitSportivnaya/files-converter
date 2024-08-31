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

type ConvertType string

const (
	SrcConvert ConvertType = "srcConvert"
	CfConvert  ConvertType = "cfConvert"
)

type OperationType string

const (
	Add    OperationType = "add"
	Delete OperationType = "delete"
	Modify OperationType = "modify"
)

type ElementOperation struct {
	ElementName string        `mapstructure:"element_name"`
	Value       string        `mapstructure:"value,omitempty"`
	Operation   OperationType `mapstructure:"operation"`
}

type FileOperation struct {
	FileName          string             `mapstructure:"file_name"`
	ElementOperations []ElementOperation `mapstructure:"element_operations"`
}

type Configuration struct {
	PlatformVersion string          `mapstructure:"platform_version"`
	Extension       string          `mapstructure:"extension"`
	InputPath       string          `mapstructure:"input_path"`
	OutputPath      string          `mapstructure:"output_path"`
	ConversionType  ConvertType     `mapstructure:"conversion_type"`
	XMLFiles        []FileOperation `mapstructure:"xml_files"`
	OtherParam      []string        `mapstructure:"other_param"`
}

type DumpInfo struct {
	ConfigName string
	Version    string
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

func LoadConfig(vp *viper.Viper) (*Configuration, error) {
	if err := vp.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Configuration
	if err := vp.Unmarshal(&config); err != nil {
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
	expandedPath := os.ExpandEnv(strings.TrimSpace(input))

	absPath, err := filepath.Abs(expandedPath)
	if err != nil {
		return "", err
	}

	cleanPath := filepath.Clean(absPath)

	return cleanPath, nil
}
