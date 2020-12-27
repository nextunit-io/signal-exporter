package signal

import (
	"io/ioutil"
	"log"
	"os"
)

type signal struct {
	path   string
	tmpDir string
	config config
}

func New(signalPath string) *signal {
	s := &signal{
		path: signalPath,
	}

	s.prepare()

	return s
}

func (s *signal) prepare() {
	s.getConfig()

	directory, err := ioutil.TempDir(os.TempDir(), "export")
	if err != nil {
		log.Fatal(err)
	}

	s.tmpDir = directory
	s.copyDatabase()
}

func (s *signal) Execute() {
	data := s.exportDatabase()

	g, err := GetGenerator(GeneratorHTML, data)
	if err != nil {
		log.Fatal(err)
	}

	g.Generate()
}

func (s *signal) Finish() {

}
