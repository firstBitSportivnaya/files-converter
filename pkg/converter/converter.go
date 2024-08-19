package converter

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/firstBitSportivnaya/files-converter/pkg/config"
	"github.com/firstBitSportivnaya/files-converter/pkg/file_modifier"

	v8 "github.com/v8platform/api"
)

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

func ConvertFromSourceFiles(cfg *config.Configuration, sourceDir, targetDir string) error {
	version := v8.WithVersion(cfg.Version)
	tmpIB, err := createTempIB()
	if err != nil {
		return fmt.Errorf("ошибка при создании базы: %w", err)
	}
	defer tmpIB.Remove()

	tmpInfoBase := tmpIB.Infobase

	// Загрузка конфигурации из файлов
	comLoadSrc := v8.LoadConfigFromFiles(cfg.InputPath)
	err = v8.Run(tmpInfoBase, comLoadSrc, version)
	if err != nil {
		return fmt.Errorf("ошибка при загрузка конфигурации из файлов: %w", err)
	}

	tmpDir := newTempDir("", "v8_src")

	defer removeDir(tmpDir)

	// получаем исходные файлы для изменений, потом переработать на копирование каталога
	comDumpConfigToFiles := v8.DumpConfigToFiles(tmpDir)
	err = v8.Run(tmpInfoBase, comDumpConfigToFiles, version)
	if err != nil {
		return fmt.Errorf("ошибка получения исходных файлов: %w", err)
	}

	file_modifier.ChangeFiles(cfg, tmpDir)

	// Загрузка конфигурации расширения из исходников
	load := v8.LoadExtensionConfigFromFiles(tmpDir, "ПроектнаяБиблиотекаПодсистем")

	err = v8.Run(tmpInfoBase, load, version)

	if err != nil {
		return fmt.Errorf("ошибка загрузки конфигурации расширения: %w", err)
	}

	// Формирование пути для сохранения файла cfe
	outPath := filepath.Join(cfg.OutputPath, "PSSL_1_0_0_2.cfe")

	// Выгрузка в cfe
	dump := v8.DumpExtensionCfg(outPath, "ПроектнаяБиблиотекаПодсистем")

	err = v8.Run(tmpInfoBase, dump, version)
	if err != nil {
		return fmt.Errorf("ошибка при выгрузке в файл .cfe: %w", err)
	}

	fmt.Printf("файл *.cfe успешно сохранен в дирректорию: %s\n", cfg.InputPath)

	return nil
}

func ConvertFromCf(cfg *config.Configuration, sourcePath, targetDir string) error {
	version := v8.WithVersion(cfg.Version)
	tmpIB, err := createTempIB()
	if err != nil {
		return fmt.Errorf("ошибка при создании базы: %w", err)
	}
	defer tmpIB.Remove()

	tmpInfoBase := tmpIB.Infobase

	comLoadCfg := v8.LoadCfg(cfg.InputPath)
	err = v8.Run(tmpInfoBase, comLoadCfg, version)
	if err != nil {
		return fmt.Errorf("ошибка при загрузка конфигурации из файла: %w", err)
	}

	tmpDir := newTempDir("", "v8_src")
	defer removeDir(tmpDir)

	comDumpConfigToFiles := v8.DumpConfigToFiles(tmpDir)
	err = v8.Run(tmpInfoBase, comDumpConfigToFiles, version)
	if err != nil {
		return fmt.Errorf("ошибка получения исходных файлов: %w", err)
	}

	file_modifier.ChangeFiles(cfg, tmpDir)

	load := v8.LoadExtensionConfigFromFiles(tmpDir, "ПроектнаяБиблиотекаПодсистем")

	err = v8.Run(tmpInfoBase, load, version)

	if err != nil {
		return fmt.Errorf("ошибка загрузки конфигурации расширения: %w", err)
	}

	outPath := filepath.Join(cfg.OutputPath, "PSSL_1_0_0_2.cfe")

	dump := v8.DumpExtensionCfg(outPath, "ПроектнаяБиблиотекаПодсистем")

	err = v8.Run(tmpInfoBase, dump, version)
	if err != nil {
		return fmt.Errorf("ошибка при выгрузке в файл .cfe: %w", err)
	}

	fmt.Printf("файл *.cfe успешно сохранен в дирректорию: %s\n", cfg.OutputPath)

	return nil
}

func newTempDir(dir, pattern string) string {
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
