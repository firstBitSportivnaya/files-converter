package converter

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/firstBitSportivnaya/files-converter/pkg/file_modifier"

	v8 "github.com/v8platform/api"
)

func NewTempDir(dir, pattern string) string {
	tempDir, err := os.MkdirTemp(dir, pattern)
	if err != nil {
		log.Fatal(err)
	}
	return tempDir
}

func removeDir(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		log.Fatal(err)
	}
}

func ConvertFromSourceFiles(sourceDir, targetDir string) error {
	// Создание временной базы данных
	version := v8.WithVersion("8.3.23")
	tmpInfoBase, err := v8.CreateTempInfobase()
	infobasePath := strings.TrimSuffix(strings.TrimPrefix(tmpInfoBase.Connect.String(), "File='"), "';")
	defer removeDir(infobasePath)

	// Обработка ошибки создания базы
	if err != nil {
		return fmt.Errorf("ошибка при создании базы: %w", err)
	}

	// Загрузка конфигурации из файлов
	comLoadSrc := v8.LoadConfigFromFiles(sourceDir)
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

	// Формирование пути для сохранения файла cfe
	outPath := filepath.Join(targetDir, "PSSL_1_0_0_2.cfe")

	// Выгрузка в cfe
	dump := v8.DumpExtensionCfg(outPath, "ПроектнаяБиблиотекаПодсистем")

	err = v8.Run(tmpInfoBase, dump, version)
	if err != nil {
		return fmt.Errorf("ошибка при выгрузке в файл .cfe: %w", err)
	}

	fmt.Printf("файл *.cfe успешно сохранен в дирректорию: %s\n", sourceDir)

	return nil
}

type TempInfoBase struct {
	Infobase     *v8.Infobase
	infobasePath string
}

func (ib *TempInfoBase) SetPath() {
	if ib.Infobase != nil {
		ib.infobasePath = strings.TrimSuffix(strings.TrimPrefix(ib.Infobase.Connect.String(), "File='"), "';")
	}
}

func (ib *TempInfoBase) GetPath() string {
	return ib.infobasePath
}

func (ib *TempInfoBase) Remove() {
	if ib.infobasePath != "" {
		removeDir(ib.infobasePath)
	}
}

func createTempIB() (*TempInfoBase, error) {
	tmpInfoBase, err := v8.CreateTempInfobase()
	if err != nil {
		return nil, err
	}

	infobase := &TempInfoBase{
		Infobase: tmpInfoBase,
	}

	infobase.SetPath()

	return infobase, nil
}

func ConvertFromCf(sourcePath, targetDir string) error {
	version := v8.WithVersion("8.3.23")
	tmpIB, err := createTempIB()
	if err != nil {
		return fmt.Errorf("ошибка при создании базы: %w", err)
	}
	defer tmpIB.Remove()

	tmpInfoBase := tmpIB.Infobase

	comLoadCfg := v8.LoadCfg(sourcePath)
	err = v8.Run(tmpInfoBase, comLoadCfg, version)
	if err != nil {
		return fmt.Errorf("ошибка при загрузка конфигурации из файла: %w", err)
	}

	tmpDir := NewTempDir("", "v8_src")
	defer removeDir(tmpDir)

	comDumpConfigToFiles := v8.DumpConfigToFiles(tmpDir)
	err = v8.Run(tmpInfoBase, comDumpConfigToFiles, version)
	if err != nil {
		return fmt.Errorf("ошибка получения исходных файлов: %w", err)
	}

	file_modifier.ChangeFiles(tmpDir)

	load := v8.LoadExtensionConfigFromFiles(tmpDir, "ПроектнаяБиблиотекаПодсистем")

	err = v8.Run(tmpInfoBase, load, version)

	if err != nil {
		return fmt.Errorf("ошибка загрузки конфигурации расширения: %w", err)
	}

	outPath := filepath.Join(targetDir, "PSSL_1_0_0_2.cfe")

	dump := v8.DumpExtensionCfg(outPath, "ПроектнаяБиблиотекаПодсистем")

	err = v8.Run(tmpInfoBase, dump, version)
	if err != nil {
		return fmt.Errorf("ошибка при выгрузке в файл .cfe: %w", err)
	}

	fmt.Printf("файл *.cfe успешно сохранен в дирректорию: %s\n", targetDir)

	return nil
}
