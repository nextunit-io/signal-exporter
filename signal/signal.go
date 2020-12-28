package signal

import (
	"io/ioutil"
	"log"
	"os"
)

type Signal struct {
	path   string
	tmpDir string
	config config
}

func New(signalPath string) *Signal {
	s := &Signal{
		path: signalPath,
	}

	s.prepare()

	return s
}

func (s *Signal) prepare() {
	s.getConfig()

	directory, err := ioutil.TempDir(os.TempDir(), "export")
	if err != nil {
		log.Fatal(err)
	}

	s.tmpDir = directory
	s.copyDatabase()
}

func (s *Signal) Execute() {
	data := s.exportDatabase()

	g, err := s.GetGenerator(GeneratorHTML, data)
	if err != nil {
		log.Fatal(err)
	}

	g.Generate()
}

func (s *Signal) Finish() {
	os.RemoveAll(s.tmpDir)
}

func (s *Signal) GetHomePath() string {
	return s.path
}
