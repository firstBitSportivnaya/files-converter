package export_format

import (
	"encoding/json"
	"fmt"

	_ "embed"
)

//go:embed export_format_versions.json
var exportFormat []byte

type ExportFormatVersions map[string]string

// LoadFormatVersions читает и анализирует файл JSON, содержащий сопоставления версий формата экспорта.
// Файл JSON должен содержать соответствие версий платформы для версий формата.
// Возвращает ошибку, если файл не может быть прочитан или содержит недопустимый JSON.
func LoadFormatVersions(filePath string) (ExportFormatVersions, error) {
	var versions ExportFormatVersions

	if err := json.Unmarshal(exportFormat, &versions); err != nil {
		return nil, fmt.Errorf("ошибка парсинга json файла: %w", err)
	}

	return versions, nil
}
