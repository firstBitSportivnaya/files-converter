package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/firstBitSportivnaya/files-converter/pkg/converter"
)

func main() {
	defer pressAnyKeyToExit()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите путь к исходным файлам: ")
	path, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Ошибка при чтении ввода: %v\n", err)
		return
	}

	path = strings.TrimSpace(path)
	path = filepath.Clean(path)

	fmt.Printf("Выбрана директория: %s\n", path)

	err = converter.ConvertToCfe(path)
	if err != nil {
		fmt.Printf("Ошибка при конвертации файлов: %v\n", err)
	}
}

func pressAnyKeyToExit() {
	fmt.Println("Press any key to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
