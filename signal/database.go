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

func (s *Signal) getTempDatabasePath() string {
	return fmt.Sprintf("%s/%s", s.tmpDir, "db.sqlite")
}

func (s *Signal) copyDatabase() {
	from, err := os.Open(s.path + "/sql/db.sqlite")
	if err != nil {
		log.Print(color.RedString("[FAILED] "), "Make safety copy of database")
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(s.getTempDatabasePath(), os.O_RDWR|os.O_CREATE, 0666)
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

func (s *Signal) exportDatabase() *models.SignalData {
	data := &models.SignalData{
		Conversations: s.exportConversations(),
	}

	log.Print(color.GreenString("[OK] "), fmt.Sprintf("Got %d conversations", len(data.Conversations)))

	messageCount := 0
	for i := range data.Conversations {
		data.Conversations[i].Messages = s.exportMessages(data.Conversations[i].ID)
		messageCount += len(data.Conversations[i].Messages)
	}

	log.Print(color.GreenString("[OK] "), fmt.Sprintf("Got %d messages", messageCount))

	log.Print(color.GreenString("[OK] "), "Export Database")

	return data
}

func (s *Signal) exportMessages(conversationID string) []models.SignalMessage {
	var b bytes.Buffer
	var database bytes.Buffer

	cmd := exec.Command("sqlcipher", "-list", "-noheader", s.getTempDatabasePath(), fmt.Sprintf("PRAGMA key = \"x'%s'\";select json from messages WHERE conversationId = '%s';", s.config.Key, conversationID))
	cmd.Stdout = &database
	cmd.Stderr = &b

	err := cmd.Run()

	if err != nil {
		log.Print(color.RedString("[FAILED] "), "Getting database input for messages")
		log.Fatal(err, b.String())
		return nil
	}

	log.Print(color.GreenString("[OK] "), "Getting database input for messages")

	data := []models.SignalMessage{}
	scanner := bufio.NewScanner(strings.NewReader(database.String()))
	var d models.SignalMessage

	for scanner.Scan() {
		if scanner.Text() == "ok" {
			continue
		}

		err := json.Unmarshal(scanner.Bytes(), &d)
		if err != nil {
			log.Print(color.RedString("[FAILED] "), "Parsing messages")
			log.Fatal("Cannot parse", err, scanner.Text())
			return nil
		}

		data = append(data, d)
	}

	return data
}

func (s *Signal) exportConversations() []models.SignalConverstation {
	var b bytes.Buffer
	var database bytes.Buffer

	cmd := exec.Command("sqlcipher", "-list", "-noheader", s.getTempDatabasePath(), "PRAGMA key = \"x'"+s.config.Key+"'\";select json from conversations;")
	cmd.Stdout = &database
	cmd.Stderr = &b

	err := cmd.Run()

	if err != nil {

		log.Print(color.RedString("[FAILED] "), "Getting database input for conversations")
		log.Fatal(err, b.String())
		return nil
	}

	log.Print(color.GreenString("[OK] "), "Getting database input for conversations")

	data := []models.SignalConverstation{}

	scanner := bufio.NewScanner(strings.NewReader(database.String()))
	var d models.SignalConverstation

	for scanner.Scan() {
		if scanner.Text() == "ok" {
			continue
		}

		err := json.Unmarshal(scanner.Bytes(), &d)
		if err != nil {
			log.Print(color.RedString("[FAILED] "), "Parsing conversations")
			log.Fatal("Cannot parse", err, scanner.Text())
			return nil
		} else {
			data = append(data, d)
		}
	}

	return data
}
