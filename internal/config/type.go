package config

const NamePrefixElement = "NamePrefix"

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
	ElementName string        `json:"element_name"`
	Value       string        `json:"value,omitempty"`
	Operation   OperationType `json:"operation"`
}

type FileOperation struct {
	FileName          string              `json:"file_name"`
	ElementOperations []*ElementOperation `json:"element_operations"`
}

type Configuration struct {
	PlatformVersion string           `json:"platform_version"`
	Extension       string           `json:"extension"`
	Prefix          string           `json:"prefix"`
	InputPath       string           `json:"input_path" env-required:"true"`
	OutputPath      string           `json:"output_path" env-required:"true"`
	ConversionType  ConvertType      `json:"conversion_type" env-required:"true"`
	XMLFiles        []*FileOperation `json:"xml_file_changes"`
}

type DumpInfo struct {
	ConfigName string
	Version    string
}
