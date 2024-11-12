package export_format

import (
	"encoding/json"
	"os"
)

type (
	platform_version string
	format_version   string
)

type ExportFomatVersions map[string]string

func LoadFormatVersions(filePath string) (ExportFomatVersions, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var versions ExportFomatVersions

	if err := json.Unmarshal(data, &versions); err != nil {
		return nil, err
	}

	return versions, nil
}
