package file_modifier

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/beevik/etree"
)

var (
	filesToProcess = []string{
		// main file
		"Configuration.xml",
		// Languages
		"Русский.xml",
		"English.xml",
		// Roles
		"АдминистраторСистемы.xml",
		"ИнтерактивноеОткрытиеВнешнихОтчетовИОбработок.xml",
		"ПолныеПрава.xml",
	}
	dirCommonModules = "CommonModules"
)

func ChangeFiles(dir string) {

	files := make(map[string]struct{})
	for _, fileName := range filesToProcess {
		files[fileName] = struct{}{}
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
				if _, found := files[dirEntryName]; found {
					return processFile(path, dirEntryName)
				}
			}
		}
		return nil
	}

	// Проходим по всем файлам и поддиректориям в корневой директории
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
			if err := DisablePrivilegedMode(filePath, entryName); err != nil {
				return err
			}
		}
	}
	return nil
}

func isXMLFile(fileName string) bool {
	return filepath.Ext(fileName) == ".xml"
}

func DisablePrivilegedMode(path, fileName string) error {
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
		if err := changeAttribute(properties, key, "false"); err != nil {
			return err
		}
		if err := doc.WriteToFile(path); err != nil {
			return fmt.Errorf("ошибка при записи файла %s: %w", path, err)
		}

		fmt.Println("Файл успешно обработан:", fileName)
	}

	return nil
}

func processFile(path, fileName string) error {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		return fmt.Errorf("ошибка при чтении файла %s: %w", fileName, err)
	}

	properties := findProperties(doc)
	if properties == nil {
		return fmt.Errorf("элемент <Properties> не найден в файле %s", fileName)
	}

	// Обработка в зависимости от типа файла
	switch filepath.Base(path) {
	case "Configuration.xml":
		if err := processConfigurationFile(properties); err != nil {
			return fmt.Errorf("ошибка при обработке Configuration.xml: %w", err)
		}
	default:
		err := addAttributes(properties, map[string]string{
			"ObjectBelonging": "Adopted",
		})
		if err != nil {
			return fmt.Errorf("ошибка при добавлении атрибутов в файл %s: %w", fileName, err)
		}
	}

	// Записываем изменения обратно в файл
	if err := doc.WriteToFile(path); err != nil {
		return fmt.Errorf("ошибка при записи файла %s: %w", path, err)
	}

	fmt.Println("Файл успешно обработан:", fileName)
	return nil
}

// findProperties - Находит элемент <Properties>
func findProperties(doc *etree.Document) *etree.Element {
	return doc.FindElement("//Properties")
}

// processConfigurationFile - Обрабатывает файл Configuration.xml
func processConfigurationFile(properties *etree.Element) error {

	var elementsToDelete = []string{
		"DefaultRoles",
		"DefaultRunMode",
		"UsePurposes",
		"DefaultReportForm",
		"DefaultReportVariantForm",
		"DefaultReportSettingsForm",
		"DefaultReportAppearanceTemplate",
		"DefaultDynamicListSettingsForm",
		"DefaultSearchForm",
		"DefaultDataHistoryChangeHistoryForm",
		"DefaultDataHistoryVersionDataForm",
		"DefaultDataHistoryVersionDifferencesForm",
		"DefaultCollaborationSystemUsersChoiceForm",
		"DefaultStyle",
		"ModalityUseMode",
		"SynchronousPlatformExtensionAndAddInCallUseMode",
		"InterfaceCompatibilityMode",
		"DatabaseTablespacesUseMode",
		"CompatibilityMode",
	}

	// Добавление и удаление атрибутов
	err := addAttributes(properties, map[string]string{
		"ObjectBelonging": "Adopted",
		"NamePrefix":      "пбп_",
	})
	if err != nil {
		return err
	}

	err = changeAttribute(properties, "ConfigurationExtensionCompatibilityMode", "Version8_3_21")
	if err != nil {
		return err
	}

	err = removeElements(properties, elementsToDelete)
	if err != nil {
		return err
	}

	return nil
}

// addAttributes - Добавляет атрибуты к элементу
func addAttributes(properties *etree.Element, attributes map[string]string) error {
	for key, value := range attributes {
		elem := properties.FindElement(key)
		if elem == nil {
			elem = properties.CreateElement(key)
		}
		elem.SetText(value)
	}
	return nil
}

// changeAttribute - Изменяет значение атрибута
func changeAttribute(properties *etree.Element, key, value string) error {
	elem := properties.FindElement(key)
	if elem == nil {
		return fmt.Errorf("атрибут %s не найден", key)
	}
	elem.SetText(value)
	return nil
}

// removeElements - Удаляет указанные элементы
func removeElements(properties *etree.Element, elements []string) error {
	for _, elem := range elements {
		element := properties.FindElement(elem)
		if element != nil {
			properties.RemoveChild(element)
		}
	}
	return nil
}
