package utils

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"strings"
)

func CreateTemlpate(tmplContent string) (tmplParse *template.Template, err error) {

	tmpl := template.New("output")

	tmplParse, err = tmpl.Parse(tmplContent)
	if err != nil {
		return nil, err
	}
	return tmplParse, nil
}

func CreateFile(filePath string) (file *os.File, err error) {

	file, err = os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func ReadFile(filePath string) (file *os.File, err error) {
	file, err = os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return nil, err
	}
	return file, nil
}

func ParseFile(file *os.File) {
	scanner := bufio.NewScanner(file)
	// Читаем файл построчно
	for scanner.Scan() {
		line := scanner.Text()
		s := strings.TrimRight(line, ":")

		fmt.Println(s)
	}
	defer file.Close()
}
