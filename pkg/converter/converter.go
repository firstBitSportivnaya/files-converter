package converter

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/firstBitSportivnaya/files-converter/pkg/config"
	"github.com/firstBitSportivnaya/files-converter/pkg/utils/fileutil"
	xmlutil "github.com/firstBitSportivnaya/files-converter/pkg/utils/xmlutil"

	v8 "github.com/v8platform/api"
	"github.com/v8platform/runner"
)

type SourceFileConverter struct{}

func (s *SourceFileConverter) Convert(cfg *config.Configuration) error {
	if err := ConvertToCfe(cfg); err != nil {
		return err
	}
	return nil
}

type CfConverter struct{}

func (c *CfConverter) Convert(cfg *config.Configuration) error {
	if err := ConvertToCfe(cfg); err != nil {
		return err
	}
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

func ConvertToCfe(cfg *config.Configuration) error {
	dumpInfo := config.GetDumpInfo()

	version := v8.WithVersion(cfg.PlatformVersion)
	tmpIB, err := createTempIB()
	if err != nil {
		return err
	}
	defer tmpIB.Remove()

	tmpInfoBase := tmpIB.Infobase

	tmpDir := newTempDir("", "v8_src")
	defer removeDir(tmpDir)

	switch cfg.ConversionType {
	case config.SrcConvert:
		if err = fileutil.CopyDir(cfg.InputPath, tmpDir); err != nil {
			return err
		}
	case config.CfConvert:
		if err = loadCfConfig(cfg, tmpInfoBase, version); err != nil {
			return err
		}

		comDumpConfigToFiles := v8.DumpConfigToFiles(tmpDir)
		if err = v8.Run(tmpInfoBase, comDumpConfigToFiles, version); err != nil {
			return fmt.Errorf("ошибка получения исходных файлов: %w", err)
		}
	}

	if err = xmlutil.ChangeFiles(cfg, tmpDir); err != nil {
		return err
	}

	extension := cfg.Extension
	if extension == "" {
		extension = dumpInfo.ConfigName
	}

	load := v8.LoadExtensionConfigFromFiles(tmpDir, extension)
	if err = v8.Run(tmpInfoBase, load, version); err != nil {
		return fmt.Errorf("ошибка загрузки конфигурации расширения: %w", err)
	}

	outputFile := extension
	if dumpInfo.Version != "" {
		outputFile += "_" + strings.ReplaceAll(dumpInfo.Version, ".", "_")
	}
	outputFile += ".cfe"
	outPath := filepath.Join(cfg.OutputPath, outputFile)

	dump := v8.DumpExtensionCfg(outPath, extension)
	if err = v8.Run(tmpInfoBase, dump, version); err != nil {
		return fmt.Errorf("ошибка при выгрузке в файл .cfe: %w", err)
	}

	fmt.Printf("файл *.cfe успешно сохранен в дирректорию: %s\n", cfg.InputPath)

	return nil
}

func loadCfConfig(cfg *config.Configuration, tmpInfoBase *v8.Infobase, version runner.Option) error {
	comLoadCfg := v8.LoadCfg(cfg.InputPath)
	if err := v8.Run(tmpInfoBase, comLoadCfg, version); err != nil {
		return fmt.Errorf("ошибка при загрузка конфигурации из файла: %w", err)
	}
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
		return nil, fmt.Errorf("ошибка при создании базы: %w", err)
	}

	infobase := &TempInfoBase{
		Infobase: tmpInfoBase,
	}

	infobase.SetPath()

	return infobase, nil
}
