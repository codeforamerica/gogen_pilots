package processor

import (
	"encoding/csv"
	"gogen/data"
	"os"
	"path/filepath"
)

type DataProcessor struct {
	outputFolder string
}

func NewDataProcessor(
	cmsCSV *csv.Reader,
	weightsInformation *data.WeightsInformation,
	dojInformation *data.DOJInformation,
	outputFolder string,
) DataProcessor {
	return DataProcessor{outputFolder: outputFolder}
}

func (d DataProcessor) Process() {
	os.Create(filepath.Join(d.outputFolder, "felonies_sf_results.csv"))
}
