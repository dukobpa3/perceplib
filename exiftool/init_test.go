package exiftool

import (
	"log"
	"runtime"
)

func init() {
	if runtime.GOOS == "windows" {
		Exec = `dist\exiftool.exe`
		Arg1 = `dist\exiftool`
	} else {
		Exec = "dist/exiftool"
	}
	log.Println("Init from exiftool")
}
