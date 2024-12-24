package export_format

import (
	"encoding/json"
	"fmt"

	_ "embed"
)

//go:embed export_format_versions.json
var exportFormat []byte

type ExportFormatVersions map[string]string

// LoadFormatVersions возвращает сопоставления версий формата экспорта из встроенных данных.
// Возвращает ошибку, если встроенные данные содержат недопустимый JSON.
func LoadFormatVersions(filePath string) (ExportFormatVersions, error) {
	var versions ExportFormatVersions

	if err := json.Unmarshal(exportFormat, &versions); err != nil {
		return nil, fmt.Errorf("ошибка парсинга json файла: %w", err)
	}

	return versions, nil
}
