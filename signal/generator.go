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

const outputPath = "./outputs"

type Generator interface {
	Generate() error
}

func (s *Signal) GetGenerator(generatorType GeneratorType, data *models.SignalData) (Generator, error) {
	switch generatorType {
	case GeneratorHTML:
		return html.New(s.path, data, outputPath), nil
	default:
		return nil, errors.New("Invalid generator type")
	}

}
