package logs

import (
	"io"
	"log"
	"os"

	"github.com/spf13/viper"
)

type LogFile struct {
	file *os.File
}

func InitLogs() *LogFile {

	f, err := os.OpenFile(viper.GetString("paths.logs"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	return &LogFile{file: f}
}

func (f *LogFile) CloseFile() {

	f.file.Close()
}
