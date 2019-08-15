package utilities

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	OTHER_ERROR                      = 1
	FILE_PARSING_ERROR               = 2
	INVALID_RUN_OPTION_ERROR         = 3
	INVALID_ELIGIBILITY_OPTION_ERROR = 4
)

type GogenError struct {
	ExitCode     int
	ErrorMessage string
}

func (g *GogenError) Error() string {
	return g.ErrorMessage
}

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
	errorMap := map[string]GogenError{
		"": {ExitCode: exitCode, ErrorMessage: originalError.Error()},
	}
	ExitWithErrors(errorMap)
}

func ExitWithErrors(originalErrors map[string]GogenError) {
	s, err := json.Marshal(originalErrors)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(OTHER_ERROR)
	}

	err = ioutil.WriteFile(errorFileName, s, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(OTHER_ERROR)
	}

	var exitCode int

	for _, gogenError := range originalErrors {
		if exitCode == 0 || gogenError.ExitCode == OTHER_ERROR{
			exitCode = gogenError.ExitCode
		}
	}

	fileNames := make([]string, 0, len(originalErrors))
	for key := range originalErrors {
		fileNames = append(fileNames, key)
	}
	sort.Strings(fileNames)

	for _, fileName := range fileNames {
		fmt.Fprintf(os.Stderr, "%s: %s\n", fileName, originalErrors[fileName].ErrorMessage)
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
