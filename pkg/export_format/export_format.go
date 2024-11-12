package export_format

import (
	"encoding/json"
	"fmt"
	"os"
)

type ExportFormatVersions map[string]string

// LoadFormatVersions читает и анализирует файл JSON, содержащий сопоставления версий формата экспорта.
// Файл JSON должен содержать соответствие версий платформы для версий формата.
// Возвращает ошибку, если файл не может быть прочитан или содержит недопустимый JSON.
func LoadFormatVersions(filePath string) (ExportFormatVersions, error) {
	if filePath == "" {
		return nil, fmt.Errorf("путь к файлу не может быть пустым")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла версий формата: %w", err)
	}

	var versions ExportFormatVersions

	if err := json.Unmarshal(data, &versions); err != nil {
		return nil, fmt.Errorf("ошибка парсинга json файла: %w", err)
	}

	return versions, nil
}
