package processor

import "gogen/data"

type DataProcessor struct{}

func NewDataProcessor(
	cmsInformation *data.CMSInformation,
	weightsInformation *data.WeightsInformation,
	dojInformation *data.DOJInformation,
	) DataProcessor {
	return DataProcessor{}
}

func (d DataProcessor) Process() {}