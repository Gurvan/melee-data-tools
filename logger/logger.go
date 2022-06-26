package logger

import (
	"log"
	"os"
)

var (
	Warning *log.Logger
	Info    *log.Logger
	Error   *log.Logger
)

func init() {
	Info = log.New(os.Stdout, "INFO: ", log.Lshortfile)
	Warning = log.New(os.Stdout, "WARNING: ", log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ", log.Lshortfile)
}
