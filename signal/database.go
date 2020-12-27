package signal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"nextunit.io/signal-exporter/models"
)

func (s *signal) copyDatabase() {
	from, err := os.Open(s.path + "/sql/db.sqlite")
	if err != nil {
		log.Print(color.RedString("[FAILED] "), "Make safety copy of database")
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(s.tmpDir+"/db.sqlite", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Print(color.RedString("[FAILED] "), "Make safety copy of database")
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Print(color.RedString("[FAILED] "), "Make safety copy of database")
		log.Fatal(err)
	}

	log.Print(color.GreenString("[OK] "), "Make safety copy of database")
}

func (s *signal) exportDatabase() *models.SignalData {
	path := s.tmpDir + "/db.sqlite"

	var b bytes.Buffer
	var database bytes.Buffer

	cmd := exec.Command("sqlcipher", "-list", "-noheader", path, "PRAGMA key = \"x'"+s.config.Key+"'\";select json from messages;")
	cmd.Stdout = &database
	cmd.Stderr = &b

	err := cmd.Run()

	if err != nil {
		log.Fatal(err, b.String())
		return nil
	}

	data := &models.SignalData{
		Conversations: make(map[string]*models.SignalConverstation),
	}

	log.Print(color.GreenString("[OK] "), "Getting database input")

	scanner := bufio.NewScanner(strings.NewReader(database.String()))
	for scanner.Scan() {
		if scanner.Text() == "ok" {
			continue
		}

		var d models.SignalMessage

		err := json.Unmarshal(scanner.Bytes(), &d)
		if err != nil {
			log.Print(color.RedString("[FAILED] "), "Parsing messages")
			log.Fatal("Cannot parse", err, scanner.Text())
			return nil
		} else {
			if _, ok := data.Conversations[d.ConversationID]; !ok {
				data.Conversations[d.ConversationID] = &models.SignalConverstation{
					ConversationID: d.ConversationID,
					Messages:       []models.SignalMessage{},
				}
			}

			data.Conversations[d.ConversationID].Messages = append(data.Conversations[d.ConversationID].Messages, d)
		}
	}

	log.Print(color.GreenString("[OK] "), fmt.Sprintf("Parsing done. Found %d", len(data.Conversations)))

	return data
}
