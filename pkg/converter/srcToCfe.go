package converter

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/firstBitSportivnaya/files-converter/pkg/file_modifier"

	v8 "github.com/v8platform/api"
)

// NewTempDir получает имя нового временного каталога
func NewTempDir(dir, pattern string) string {
	tempDir, err := os.MkdirTemp(dir, pattern)
	if err != nil {
		log.Fatal(err)
	}
	return tempDir
}

func removeDir(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		log.Fatal(err)
	}
}

func ConvertToCfe(dir string) error {
	reader := bufio.NewReader(os.Stdin)

	// Создание временной базы данных
	version := v8.WithVersion("8.3.23")
	tmpInfoBase, err := v8.CreateTempInfobase()
	// Получение пути к временной базе данных
	infobasePath := strings.Replace(tmpInfoBase.Connect.String(), "File='", "", -1)
	infobasePath = strings.Replace(infobasePath, "';", "", -1)

	// Удаление временной базы данных
	defer removeDir(infobasePath)

	// Обработка ошибки создания базы
	if err != nil {
		return fmt.Errorf("ошибка при создании базы: %w", err)
	}

	// Загрузка конфигурации из файлов
	comLoadSrc := v8.LoadConfigFromFiles(dir)
	err = v8.Run(tmpInfoBase, comLoadSrc, version)
	if err != nil {
		return fmt.Errorf("ошибка при загрузка конфигурации из файлов: %w", err)
	}

	// временный каталог для исходных файлов
	tmpDir := NewTempDir("", "v8_src")

	// Удаление временной папки с исходными файлами
	defer removeDir(tmpDir)

	// получаем исходные файлы для изменений, потом переработать на копирование каталога
	comDumpConfigToFiles := v8.DumpConfigToFiles(tmpDir)
	err = v8.Run(tmpInfoBase, comDumpConfigToFiles, version)
	if err != nil {
		return fmt.Errorf("ошибка получения исходных файлов: %w", err)
	}

	// Обработка файлов
	file_modifier.ChangeFiles(tmpDir)

	// Загрузка конфигурации расширения из исходников
	load := v8.LoadExtensionConfigFromFiles(tmpDir, "ПроектнаяБиблиотекаПодсистем")

	err = v8.Run(tmpInfoBase, load, version)

	if err != nil {
		return fmt.Errorf("ошибка загрузки конфигурации расширения: %w", err)
	}

	// Ввод пути для сохранения файла cfe
	fmt.Print("Введите путь для сохранения файла *cfe: ")
	savePath, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка при чтении ввода: %w", err)
	}
	savePath = strings.TrimSpace(savePath)
	savePath = filepath.Clean(savePath)

	// Формирование пути для сохранения файла cfe
	dirOut := filepath.Join(savePath, "PSSL_1_0_0_2.cfe")

	// Выгрузка в cfe
	dump := v8.DumpExtensionCfg(dirOut, "ПроектнаяБиблиотекаПодсистем")

	err = v8.Run(tmpInfoBase, dump, version)
	if err != nil {
		return fmt.Errorf("ошибка при выгрузке в файл .cfe: %w", err)
	}

	fmt.Printf("файл *cfe успешно сохранен в дирректорию: %s\n", dir)

	return nil
}
