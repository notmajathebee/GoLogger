/*
	---------------------------
	Title: Logger Package

	Author: MHenninger
	Version: 1.1.2
	Created: 2025-03-13
	Updated: 2025-03-20
	---------------------------
			User Guide:

	// Innerhalb main() fill Handler structur
	file, err := logger.InitializeLogger()
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Aufruf der LogHandler:
	- logger.Info()
	- logger.Warning()
	- logger.Error()

	- logger.PanicHandler()
	// main() -> defer logger.PanicHandler()
	// Damit Panics auch in der log file angezeigt werden
*/

package loggerr

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	InfoLog    *log.Logger
	WarningLog *log.Logger
	ErrorLog   *log.Logger
)

func InitializeLogger() (*os.File, error) {

	// Datum formatieren
	date := time.Now().Format("2006-01-02")

	// Get dir path
	dirPath, _ := os.Getwd()
	dirPath = filepath.Join(dirPath, "logfiles")

	// Überprüft ob folder existiert oder erstellt werden muss
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// File zusammensetzen (Namen bestimmen)
	fileName := fmt.Sprintf("%s/%s.log", dirPath, date)

	// Log in eine Datei schreiben
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		//log.Fatal(err)
		return nil, err
	} //defer file.Close()

	log.SetOutput(file)

	// Create multiWriter (Write to two files)
	multiWriter := io.MultiWriter(os.Stdout, file)

	// Config LogHandler
	InfoLog = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lmsgprefix)
	WarningLog = log.New(multiWriter, "WARNING: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lmsgprefix)
	ErrorLog = log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lmsgprefix)

	return file, nil
}

// Config functions from LogHandler
func Info(a ...any) {
	InfoLog.Println(a...)
}
func Warning(a ...any) {
	WarningLog.Println(a...)
}
func Error(a ...any) {
	// (1) Hollt den vorletzten "Hop" aus dem Stack("History") (0) Letzten Hop (logger.go)
	_, secfile, secline, ok2 := runtime.Caller(2)
	_, firfile, firline, ok1 := runtime.Caller(1)
	if ok1 && ok2 {
		shortFirFile := filepath.Base(firfile)
		shortSecFile := filepath.Base(secfile)
		ErrorLog.Printf("%s:%d -> %s:%d: %v", shortSecFile, secline, shortFirFile, firline, fmt.Sprint(a...))
	} else {
		ErrorLog.Println(a...)
	}

}

// Panic-Handler: Fängt Panics am Ende ab und schreibt sie ins Log
func PanicHandler() {
	if r := recover(); r != nil {
		ErrorLog.Println("Panic:", r)
	}
} //main() -> defer logger.HandlePanic()
