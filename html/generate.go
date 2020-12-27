package html

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/fatih/color"
)

const templatePath = "html/template.html"

func (data *html) Generate() error {
	os.MkdirAll(fmt.Sprintf("%s/html", data.outputPath), os.ModePerm)

	for _, c := range data.data.Conversations {
		t := template.Must(template.ParseFiles(templatePath))

		f, err := os.Create(fmt.Sprintf("%s/html/%s.html", data.outputPath, c.ID))
		if err != nil {
			log.Println("create file: ", err)
			return err
		}

		err = t.Execute(f, c)
		if err != nil {
			log.Print("execute: ", err)
			return err
		}
	}

	log.Print(color.GreenString("[OK] "), fmt.Sprintf("Generated %d files", len(data.data.Conversations)))
	return nil
}
