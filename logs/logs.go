package logs

import (
	"io"
	"log"
	"os"
)

type LogFile struct {
	file *os.File
}

func InitLogs(path string) *LogFile {

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	// return &LogFile{file: f}
	return &LogFile{file: f}
}

func (f *LogFile) CloseFile() {

	f.file.Close()
}
