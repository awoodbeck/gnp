package ch13

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func Example_log() {
	l := log.New(os.Stdout, "example: ", log.Lshortfile)
	l.Print("logging to standard output")

	// Output:
	// example: log_test.go:12: logging to standard output
}

func Example_logMultiWriter() {
	logFile := new(bytes.Buffer)
	w := SustainedMultiWriter(os.Stdout, logFile)
	l := log.New(w, "example: ", log.Lshortfile|log.Lmsgprefix)

	fmt.Println("standard output:")
	l.Print("Canada is south of Detroit")

	fmt.Print("\nlog file contents:\n", logFile.String())

	// Output:
	// standard output:
	// log_test.go:24: example: Canada is south of Detroit
	//
	// log file contents:
	// log_test.go:24: example: Canada is south of Detroit
}

func Example_logLevels() {
	lDebug := log.New(os.Stdout, "DEBUG: ", log.Lshortfile)
	logFile := new(bytes.Buffer)
	w := SustainedMultiWriter(logFile, lDebug.Writer())
	lError := log.New(w, "ERROR: ", log.Lshortfile)

	fmt.Println("standard output:")
	lError.Print("cannot communicate with the database")
	lDebug.Print("you cannot hum while holding your nose")

	fmt.Print("\nlog file contents:\n", logFile.String())

	// Output:
	// standard output:
	// ERROR: log_test.go:43: cannot communicate with the database
	// DEBUG: log_test.go:44: you cannot hum while holding your nose
	//
	// log file contents:
	// ERROR: log_test.go:43: cannot communicate with the database
}
