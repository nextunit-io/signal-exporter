package html

import (
	"nextunit.io/signal-exporter/models"
)

type html struct {
	outputPath string `default:"outputs/"`
	data       []models.SignalData
}

func New(signalData []models.SignalData) *html {
	return &html{
		data: signalData,
	}
}
