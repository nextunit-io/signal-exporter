package html

import (
	"nextunit.io/signal-exporter/models"
)

type html struct {
	outputPath string
	signalPath string
	data       *models.SignalData
}

func New(signalPath string, signalData *models.SignalData, outputPath string) *html {
	return &html{
		signalPath: signalPath,
		data:       signalData,
		outputPath: outputPath,
	}
}
