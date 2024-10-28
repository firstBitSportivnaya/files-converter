package config

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
	FileName          string              `mapstructure:"file_name"`
	ElementOperations []*ElementOperation `mapstructure:"element_operations"`
}

type Configuration struct {
	PlatformVersion string           `mapstructure:"platform_version"`
	Extension       string           `mapstructure:"extension"`
	Prefix          string           `mapstructure:"prefix"`
	InputPath       string           `mapstructure:"input_path"`
	OutputPath      string           `mapstructure:"output_path"`
	ConversionType  ConvertType      `mapstructure:"conversion_type"`
	XMLFiles        []*FileOperation `mapstructure:"xml_file_changes"`
}

type DumpInfo struct {
	ConfigName string
	Version    string
}
