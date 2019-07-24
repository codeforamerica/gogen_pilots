package utilities

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"
)

const (
	OTHER_ERROR = 1
	CSV_PARSING_ERROR = 2
	INVALID_OPTION_ERROR = 3
)

var fileNameSuffix string

func PrintProgressBar(index, totalRows int, totalTime time.Duration, tail string) {
	progress := float64(index) / float64(totalRows)
	bar := strings.Repeat("=", int(math.Round(progress*50.0)))
	space := strings.Repeat(" ", int(math.Round((1-progress)*50)))
	averageTime := AverageTime(totalTime, index)
	fmt.Printf("["+bar+space+"] %d/%d (avg time: %s) "+tail, int(index), int(totalRows), averageTime)
	fmt.Print("\r")
}

func AverageTime(totalTime time.Duration, index int) time.Duration {
	return time.Duration(float64(totalTime) / float64(index))
}

func Percent(num int, denom int) int {
	if denom == 0 {
		return 0
	}
	return num * 100 / denom
}

func SetFileNameSuffix(suffix string) {
	fileNameSuffix = suffix
}

func ExitWithError(originalError error, exitCode int) {
	var errorFileName string
	if fileNameSuffix != "" {
		errorFileName = fmt.Sprintf("gogen_%s.err", fileNameSuffix)
	} else {
		errorFileName = fmt.Sprintf("gogen.err")
	}
	errorFile, err := os.Create(errorFileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	errorWriter := io.MultiWriter(os.Stderr, errorFile)
	fmt.Fprintln(errorWriter, originalError)
	os.Exit(exitCode)
}
