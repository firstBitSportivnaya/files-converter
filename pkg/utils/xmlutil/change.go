package xmlutils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/beevik/etree"
	"github.com/firstBitSportivnaya/files-converter/pkg/config"
)

const (
	dirCommonModules = "CommonModules"
	mainFile         = "Configuration.xml"
	configDumpInfo   = "ConfigDumpInfo.xml"
)

type FileProcessingContext struct {
	Doc      *etree.Document
	Path     string
	FileName string
}

func ChangeFiles(cfg *config.Configuration, dir string) error {
	files := make(map[string][]*config.ElementOperation, len(cfg.XMLFiles))
	for _, file := range cfg.XMLFiles {
		files[file.FileName] = file.ElementOperations
	}

	defaultOps := defaultOperations()

	processFile := func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("ошибка при обработке файла %s: %w", path, err)
		}
		name := d.Name()

		if !d.IsDir() && isXMLFile(name) {
			doc, err := readXMLFile(path)
			if err != nil {
				return err
			}
			ctx := &FileProcessingContext{
				Doc:      doc,
				Path:     path,
				FileName: name,
			}

			if operations, found := files[name]; found {
				return processFile(ctx, operations)
			} else if isProcessableFile(dir, path, cfg.Prefix) && name != configDumpInfo {
				return processFile(ctx, defaultOps)
			} else if filepath.Base(filepath.Dir(path)) == dirCommonModules {
				return processCommonModules(ctx)
			} else if name == mainFile {
				return getInfoFromMainFile(ctx)
			}
		}
		return nil
	}

	if err := filepath.WalkDir(dir, processFile); err != nil {
		return fmt.Errorf("ошибка при обходе директорий: %v", err)
	}

	return nil
}

func GetFormatVersion(path string) (string, error) {
	doc, err := readXMLFile(filepath.Join(path, mainFile))
	if err != nil {
		return "", err
	}

	metaDataObject := doc.SelectElement("MetaDataObject")

	return metaDataObject.SelectAttrValue("version", ""), nil
}

func defaultOperations() []*config.ElementOperation {
	operations := make([]*config.ElementOperation, 0, 1)

	operations = append(operations, config.NewElementOperation("ObjectBelonging", "Adopted", config.Add))

	return operations
}

func isProcessableFile(root, path, prefix string) bool {
	relPath, err := filepath.Rel(root, path)
	if err != nil {
		log.Printf("ошибка при получении относительного пути: %v", err)
		return false
	}
	parts := strings.Split(filepath.ToSlash(relPath), "/")

	if len(parts) > 2 {
		return false
	}

	return !containsPrefix(relPath, prefix)
}

func containsPrefix(path, prefix string) bool {
	parts := strings.Split(filepath.ToSlash(path), "/")

	for i := len(parts) - 1; i >= 0; i-- {
		if strings.HasPrefix(parts[i], prefix) || strings.HasPrefix(parts[i], "_") {
			return true
		}
	}

	return false
}

func getInfoFromMainFile(ctx *FileProcessingContext) error {
	properties := findProperties(ctx.Doc)
	if properties == nil {
		return fmt.Errorf("элемент <Properties> не найден в файле %s", ctx.FileName)
	}
	getInfo(properties)
	return nil
}

func getInfo(properties *etree.Element) {
	dumpInfo := config.GetDumpInfo()

	currentElem := properties.FindElement("Name")
	if currentElem != nil {
		dumpInfo.SetConfigName(currentElem.Text())
	}
	currentElem = properties.FindElement("Version")
	if currentElem != nil {
		dumpInfo.SetVersion(currentElem.Text())
	}
}

func processFile(ctx *FileProcessingContext, operations []*config.ElementOperation) error {
	properties := findProperties(ctx.Doc)
	if properties == nil {
		return fmt.Errorf("элемент <Properties> не найден в файле %s", ctx.FileName)
	}

	if ctx.FileName == mainFile {
		getInfo(properties)
	}

	for _, element := range operations {
		processElement(properties, element)
	}

	if filepath.Base(filepath.Dir(ctx.Path)) == dirCommonModules {
		if _, err := disablePrivilegedMode(properties); err != nil {
			return fmt.Errorf("ошибка при изменение привелигированного режима в файле %s: %w", ctx.FileName, err)
		}
	}

	if err := ctx.Doc.WriteToFile(ctx.Path); err != nil {
		return fmt.Errorf("ошибка при записи файла %s: %w", ctx.Path, err)
	}

	fmt.Println("Файл успешно обработан:", ctx.FileName)
	return nil
}

func processCommonModules(ctx *FileProcessingContext) error {
	properties := findProperties(ctx.Doc)
	if properties == nil {
		return fmt.Errorf("элемент <Properties> не найден в файле %s", ctx.FileName)
	}

	changed, err := disablePrivilegedMode(properties)
	if err != nil {
		return fmt.Errorf("ошибка при изменение привелигированного режима в файле %s: %w", ctx.FileName, err)
	}
	if changed {
		if err := ctx.Doc.WriteToFile(ctx.Path); err != nil {
			return fmt.Errorf("ошибка при записи файла %s: %w", ctx.Path, err)
		}
		fmt.Println("Файл успешно обработан:", ctx.FileName)
	}

	return nil
}

func disablePrivilegedMode(properties *etree.Element) (bool, error) {
	key := "Privileged"
	flag := properties.FindElement(key).Text()
	value, err := strconv.ParseBool(flag)
	if err != nil {
		return false, err
	}
	if value {
		modifyElement(properties, key, "false")
		return true, nil
	}

	return false, nil
}

// findProperties - Находит элемент <Properties>
func findProperties(doc *etree.Document) *etree.Element {
	return doc.FindElement("//Properties")
}

func processElement(properties *etree.Element, element *config.ElementOperation) {
	switch element.Operation {
	case config.Add:
		addElement(properties, element.ElementName, element.Value)
	case config.Modify:
		modifyElement(properties, element.ElementName, element.Value)
	case config.Delete:
		deleteElement(properties, element.ElementName)
	default:
		log.Printf("Неизвестная операция: %v для элемента: %s", element.Operation, element.ElementName)
	}
}

func addElement(properties *etree.Element, tag, value string) {
	if currentElem := properties.FindElement(tag); currentElem == nil {
		currentElem := properties.CreateElement(tag)
		currentElem.SetText(value)
	} else {
		modifyElement(properties, tag, value)
	}
}

func modifyElement(properties *etree.Element, path, value string) {
	currentElem := properties.FindElement(path)
	if currentElem != nil {
		currentElem.SetText(value)
	} else {
		log.Printf("Элемент %s не найден", path)
	}
}

func deleteElement(properties *etree.Element, path string) {
	currentElem := properties.FindElement(path)
	if currentElem != nil {
		properties.RemoveChild(currentElem)
	}
}

func readXMLFile(path string) (*etree.Document, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		return nil, fmt.Errorf("ошибка при чтении файла %s: %w", path, err)
	}
	return doc, nil
}

func isXMLFile(fileName string) bool {
	return filepath.Ext(fileName) == ".xml"
}
