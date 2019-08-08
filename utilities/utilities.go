package utilities

import (
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	OTHER_ERROR                      = 1
	FILE_PROCESSING_ERROR            = 2
	INVALID_RUN_OPTION_ERROR         = 3
	INVALID_ELIGIBILITY_OPTION_ERROR = 4
)

var errorFileName string

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

func AddMaps(map1 map[string]int, map2 map[string]int) map[string]int {
	if map1 == nil {
		map1 = make(map[string]int)
	}

	for key := range map2 {
		map1[key] = map1[key] + map2[key]
	}
	return map1
}

func SetErrorFileName(filename string) {
	errorFileName = filename
}

func ExitWithError(originalError error, exitCode int) {
	errorFile, err := os.Create(errorFileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	errorWriter := io.MultiWriter(os.Stderr, errorFile)
	fmt.Fprintln(errorWriter, originalError)
	os.Exit(exitCode)
}

func ExitWithErrors(originalErrors []error, exitCode int) {
	errorFile, err := os.Create(errorFileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	errorWriter := io.MultiWriter(os.Stderr, errorFile)
	for _, errorMessage := range originalErrors {
		fmt.Fprintln(errorWriter, errorMessage)
	}
	os.Exit(exitCode)
}

func GenerateFileName(outputFolder string, template string, suffix string) string {
	if suffix != "" {
		suffix = "_" + suffix
	}
	return filepath.Join(outputFolder, fmt.Sprintf(template, suffix))
}

func GenerateIndexedFileName(outputFolder string, template string, fileIndex int, suffix string) string {
	if suffix != "" {
		suffix = "_" + suffix
	}
	return filepath.Join(outputFolder, fmt.Sprintf(template, fileIndex, suffix))
}

func GenerateIndexedOutputFolder(outputFolder string, fileIndex int, suffix string) string {
	if suffix != "" {
		suffix = "_" + suffix
	}

	return filepath.Join(outputFolder, fmt.Sprintf("DOJ_Input_File_%d_Results%s", fileIndex, suffix))
}

func GetOutputWriter(filePath string) io.Writer {
	summaryFile, err := os.Create(filePath)
	if err != nil {
		ExitWithError(err, OTHER_ERROR)
	}
	summaryWriter := io.MultiWriter(os.Stdout, summaryFile)
	return summaryWriter
}
