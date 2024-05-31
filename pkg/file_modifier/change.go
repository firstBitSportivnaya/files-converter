package file_modifier

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/beevik/etree"
)

var filesToProcess = []string{
	// main file
	"Configuration.xml",
	// Common modules
	"пбп_ОбщегоНазначенияПолныеПрава.xml",
	// Languages
	"Русский.xml",
	"English.xml",
	// Roles
	"АдминистраторСистемы.xml",
	"ИнтерактивноеОткрытиеВнешнихОтчетовИОбработок.xml",
	"ПолныеПрава.xml",
}

func ChangeFiles(dir string) {

	files := make(map[string]struct{})
	for _, fileName := range filesToProcess {
		files[fileName] = struct{}{}
	}

	processFile := func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(d.Name()) == ".xml" {
			fileName := filepath.Base(path)
			if _, found := files[fileName]; found {
				fmt.Printf("Найден файл: %s\n", fileName)
				changeFile(path)
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

func changeFile(path string) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		log.Fatal(err)
	}

	properties := findProperties(doc)
	if properties == nil {
		fmt.Printf("Элемент <Properties> не найден в файле %s", path)
		return
	}

	// Обработка в зависимости от типа файла
	switch filepath.Base(path) {
	case "Configuration.xml":
		err := processConfigurationFile(properties)
		if err != nil {
			log.Fatal(err)
		}
	case "пбп_ОбщегоНазначенияПолныеПрава.xml":
		err := processFileToModify(properties)
		if err != nil {
			log.Fatal(err)
		}
	default:
		err := addAttributes(properties, map[string]string{
			"ObjectBelonging": "Adopted",
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	// Записываем изменения обратно в файл
	if err := doc.WriteToFile(path); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Файл успешно обработан:", path)
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

// processFileToModify - Обрабатывает файл, где необходимо изменить атрибуты
func processFileToModify(properties *etree.Element) error {
	return changeAttribute(properties, "Privileged", "false")
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
		fmt.Printf("Атрибут %s не найден", key)
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
