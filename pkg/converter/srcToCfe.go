package converter

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"project/pkg/file_modifier"
	"strings"

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

func ConvertToCfe(dir string) {
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
		log.Fatal(err)
	}

	// Загрузка конфигурации из файлов
	comLoadSrc := v8.LoadConfigFromFiles(dir)
	err = v8.Run(tmpInfoBase, comLoadSrc, version)
	if err != nil {
		log.Fatal(err)
	}

	// временный каталог для исходных файлов
	tmpDir := NewTempDir("", "v8_src")

	// Удаление временной папки с исходными файлами
	defer removeDir(tmpDir)

	// получаем исходные файлы для изменений, потом переработать на копирование каталога
	comDumpConfigToFiles := v8.DumpConfigToFiles(tmpDir)
	err = v8.Run(tmpInfoBase, comDumpConfigToFiles, version)
	if err != nil {
		log.Fatal(err)
	}

	// Обработка файлов
	file_modifier.ChangeFiles(tmpDir)

	// Загрузка конфигурации расширения из исходников
	load := v8.LoadExtensionConfigFromFiles(tmpDir, "ПроектнаяБиблиотекаПодсистем")

	err = v8.Run(tmpInfoBase, load, version)

	if err != nil {
		log.Fatal(err)
	}

	// Ввод пути для сохранения файла cfe
	fmt.Print("Введите путь для сохранения файла *cfe: ")
	savePath, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Ошибка при чтении ввода: %v\n", err)
		return
	}
	savePath = strings.TrimSpace(savePath)
	savePath = filepath.Clean(savePath)

	// Формирование пути для сохранения файла cfe
	dirOut := filepath.Join(savePath, "PSSL.cfe")

	// Выгрузка в cfe
	dump := v8.DumpExtensionCfg(dirOut, "ПроектнаяБиблиотекаПодсистем")

	err = v8.Run(tmpInfoBase, dump, version)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("файл *cfe успешно сохранен в дирректорию: %s\n", dir)
}
