package file_modifier

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/beevik/etree"
	"github.com/firstBitSportivnaya/files-converter/pkg/config"
)

var dirCommonModules = "CommonModules"

func ChangeFiles(cfg *config.Configuration, dir string) {
	files := make(map[string][]config.ElementOperation)
	for _, file := range cfg.XMLFiles {
		files[file.FileName] = file.ElementOperations
	}

	processFile := func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("ошибка при обработке файла %s: %w", path, err)
		}
		dirEntryName := d.Name()

		if d.IsDir() {
			if dirEntryName == dirCommonModules {
				return processCommonModules(path)
			}
		} else {
			if isXMLFile(dirEntryName) {
				if operations, found := files[dirEntryName]; found {
					return processFile(path, dirEntryName, operations)
				}
			}
		}
		return nil
	}

	err := filepath.WalkDir(dir, processFile)
	if err != nil {
		fmt.Printf("Ошибка при обходе директорий: %v\n", err)
	}
}

func processCommonModules(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("ошибка при чтении директории %s: %w", path, err)
	}
	for _, entry := range entries {
		entryName := entry.Name()
		if !entry.IsDir() && isXMLFile(entryName) {
			filePath := filepath.Join(path, entryName)
			if err := disablePrivilegedMode(filePath, entryName); err != nil {
				return err
			}
		}
	}
	return nil
}

func processFile(path, fileName string, operations []config.ElementOperation) error {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		return fmt.Errorf("ошибка при чтении файла %s: %w", fileName, err)
	}

	properties := findProperties(doc)
	if properties == nil {
		return fmt.Errorf("элемент <Properties> не найден в файле %s", fileName)
	}

	for _, element := range operations {
		processElement(properties, element)
	}

	if err := doc.WriteToFile(path); err != nil {
		return fmt.Errorf("ошибка при записи файла %s: %w", path, err)
	}

	fmt.Println("Файл успешно обработан:", fileName)
	return nil
}

func disablePrivilegedMode(path, fileName string) error {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		return fmt.Errorf("ошибка при чтении файла %s: %w", fileName, err)
	}

	properties := findProperties(doc)
	if properties == nil {
		return fmt.Errorf("элемент <Properties> не найден в файле %s", fileName)
	}

	key := "Privileged"
	flag := properties.FindElement(key).Text()
	value, err := strconv.ParseBool(flag)
	if err != nil {
		return fmt.Errorf("ошибка при парсинге значения флага %s: %w", flag, err)
	}
	if value {
		modifyElement(properties, key, "false")
		if err := doc.WriteToFile(path); err != nil {
			return fmt.Errorf("ошибка при записи файла %s: %w", path, err)
		}

		fmt.Println("Файл успешно обработан:", fileName)
	}

	return nil
}

// findProperties - Находит элемент <Properties>
func findProperties(doc *etree.Document) *etree.Element {
	return doc.FindElement("//Properties")
}

func processElement(properties *etree.Element, element config.ElementOperation) {
	switch element.Operation {
	case config.Add:
		if currentElem := properties.FindElement(element.ElementName); currentElem == nil {
			addElement(properties, element.ElementName, element.Value)
		} else {
			modifyElement(properties, element.ElementName, element.Value)
		}
	case config.Modify:
		modifyElement(properties, element.ElementName, element.Value)
	case config.Delete:
		deleteElement(properties, element.ElementName)
	default:
		fmt.Printf("Неизвестная операция: %v для элемента: %s", element.Operation, element.ElementName)
	}
}

func addElement(properties *etree.Element, tag, value string) {
	currentElem := properties.CreateElement(tag)
	currentElem.SetText(value)
}

func modifyElement(properties *etree.Element, path, value string) {
	currentElem := properties.FindElement(path)
	if currentElem != nil {
		currentElem.SetText(value)
	} else {
		fmt.Printf("Элемент %s не найден", path)
	}
}

func deleteElement(properties *etree.Element, path string) {
	currentElem := properties.FindElement(path)
	if currentElem != nil {
		properties.RemoveChild(currentElem)
	}
}

func isXMLFile(fileName string) bool {
	return filepath.Ext(fileName) == ".xml"
}
