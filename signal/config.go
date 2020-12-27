package signal

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/fatih/color"
)

type config struct {
	Key             string `json:"key"`
	MediaPermission bool   `json:"mediaPermission"`
}

func (s *signal) getConfig() {
	jsonFile, err := os.Open(s.path + "/config.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var c config

	json.Unmarshal(byteValue, &c)
	s.config = c

	log.Print(color.GreenString("[OK] "), "Getting configuration")
}
