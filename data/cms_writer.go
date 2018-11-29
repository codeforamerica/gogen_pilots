package data

import (
	"encoding/csv"
	"os"
	"path/filepath"
)

type CMSWriter interface {
	WriteEntry(entry CMSEntry)
	Flush()
}

type csvWriter struct {
	outputFileWriter *csv.Writer
}

func NewCMSWriter(outputFilePath string) CMSWriter {
	outputFile, err := os.Create(filepath.Join(outputFilePath))
	if err != nil {
		panic(err)
	}

	return csvWriter{
		outputFileWriter: csv.NewWriter(outputFile),
	}
}

func (cw csvWriter) WriteEntry(entry CMSEntry) {
	cw.outputFileWriter.Write(entry.RawRow)
}

func (cw csvWriter) Flush() {
	cw.outputFileWriter.Flush()
}
