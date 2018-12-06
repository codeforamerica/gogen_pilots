package utilities

import (
	"fmt"
	"math"
	"strings"
	"time"
)

//def print_progress_bar(index, total_rows)
//progress = index.to_f / total_rows
//print '[' + ('*' * (progress * 50)) + (' ' * ((1 - progress) * 50)) + ']'
//print "\r"
//end

func PrintProgressBar(index, totalRows float64, totalTime time.Duration, tail string) {
	progress := index / totalRows
	bar := strings.Repeat("=", int(math.Round(progress*50.0)))
	space := strings.Repeat(" ", int(math.Round((1-progress)*50)))
	averageTime := AverageTime(totalTime, index)
	fmt.Printf("["+bar+space+"] %d/%d (avg time: %s) "+tail, int(index), int(totalRows), averageTime)
	fmt.Print("\r")
}

func AverageTime(totalTime time.Duration, index float64) time.Duration {
	return time.Duration(float64(totalTime) / index)
}
