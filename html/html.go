package html

import (
	"nextunit.io/signal-exporter/models"
)

type html struct {
	outputPath string
	data       *models.SignalData
}

func New(signalData *models.SignalData, outputPath string) *html {
	return &html{
		data:       signalData,
		outputPath: outputPath,
	}
}
