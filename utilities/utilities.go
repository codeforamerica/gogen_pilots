package utilities

import (
	"fmt"
	"math"
	"strings"
	"time"
)

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
