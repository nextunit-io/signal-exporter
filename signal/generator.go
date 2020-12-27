package signal

import (
	"errors"

	"nextunit.io/signal-exporter/html"
	"nextunit.io/signal-exporter/models"
)

type GeneratorType int

const (
	GeneratorHTML GeneratorType = iota
)

type Generator interface {
	Generate() error
}

func GetGenerator(generatorType GeneratorType, data []models.SignalData) (Generator, error) {
	switch generatorType {
	case GeneratorHTML:
		return html.New(data), nil
	default:
		return nil, errors.New("Invalid generator type")
	}

}