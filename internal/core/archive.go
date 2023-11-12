package core

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	homeDir   = "/Users/kirill/Documents/testenv"
	zipFile   = "/Users/kirill/Documents/testenv/archive.zip"
	sourceDir = "home_dir"
)

func archiveProfiles() {
	err := os.Chdir(homeDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Создаем новый архив
	archive, err := os.Create(zipFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer archive.Close()

	// Создаем новый архиватор
	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	// Функция для рекурсивного обхода всех файлов и папок в исходной директории
	_ = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Создаем новую структуру для файла в архиве
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			fmt.Println(err)
			return err
		}

		header.Name = path

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				fmt.Println(err)
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}

		return nil
	})

	fmt.Println("Архивирование завершено:", zipFile)
}
