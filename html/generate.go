package html

import (
	"fmt"
	"html/template"
	"log"
	"os"
)

const templatePath = "html/template.html"

func (data *html) Generate() error {
	os.MkdirAll(fmt.Sprintf("%s/html", data.outputPath), os.ModePerm)

	for _, d := range data.data.Conversations {
		t := template.Must(template.ParseFiles(templatePath))

		f, err := os.Create(fmt.Sprintf("%s/html/%s.html", data.outputPath, d.ConversationID))
		if err != nil {
			log.Println("create file: ", err)
			return err
		}

		err = t.Execute(f, d)
		if err != nil {
			log.Print("execute: ", err)
			return err
		}
	}

	return nil
}
