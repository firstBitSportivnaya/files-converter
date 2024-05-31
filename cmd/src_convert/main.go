package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"project/pkg/converter"
	"strings"
)

func main() {
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

	converter.ConvertToCfe(path)
}
