package html

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"nextunit.io/signal-exporter/models"
)

const templatePath = "html/template.html"

type htmlAttachment struct {
	Attachment        models.SignalAttachment
	HTMLPath          string
	HTMLThumbnailPath string
}
type htmlMessage struct {
	Message     models.SignalMessage
	Attachments []htmlAttachment
}

type SortMessages []htmlMessage

func (message SortMessages) Len() int      { return len(message) }
func (message SortMessages) Swap(i, j int) { message[i], message[j] = message[j], message[i] }
func (message SortMessages) Less(i, j int) bool {
	return message[i].Message.Timestamp < message[j].Message.Timestamp
}

func (data *html) Generate() error {
	for _, c := range data.data.Conversations {
		functions := template.FuncMap{
			"timestamp": func(timestamp int64) string {
				t := time.Unix(timestamp/1000, 0)

				return t.Format("15:04:05")
			},
			"attachmentThumbnail": func(attachment htmlAttachment) string {
				var path string

				if attachment.HTMLThumbnailPath != "" {
					path = attachment.HTMLThumbnailPath
				} else if attachment.HTMLPath != "" {
					path = attachment.HTMLPath
				}

				return fmt.Sprintf("./attachments/%s", filepath.Base(path))
			},
		}

		t := template.Must(template.New("template.html").Funcs(functions).ParseFiles(templatePath))

		os.MkdirAll(data.getHTMLPath(c.ID), os.ModePerm)

		messages := data.prepareMessages(c)

		for dateString, message := range messages {
			sort.Sort(SortMessages(message))

			f, err := os.Create(fmt.Sprintf("%s/%s.html", data.getHTMLPath(c.ID), dateString))
			if err != nil {
				log.Println("create file: ", err)
				return err
			}

			err = t.Execute(f, struct {
				Conversation models.SignalConverstation
				Messages     []htmlMessage
			}{
				Conversation: c,
				Messages:     message,
			})

			if err != nil {
				log.Print("execute: ", err)
				return err
			}
		}
	}

	log.Print(color.GreenString("[OK] "), fmt.Sprintf("Generated %d directories", len(data.data.Conversations)))
	return nil
}

func (data *html) getHTMLPath(conversationID string) string {
	return fmt.Sprintf("%s/html/%s", data.outputPath, conversationID)
}

func (data *html) getHTMLAttachmentPath(conversationID string) string {
	return fmt.Sprintf("%s/attachments", data.getHTMLPath(conversationID))
}

func (data *html) prepareMessages(conversation models.SignalConverstation) map[string][]htmlMessage {
	d := map[string][]htmlMessage{}

	for _, message := range conversation.Messages {
		t := time.Unix(message.Timestamp/1000, 0)

		timeString := t.Format("2006-01-02")

		if _, ok := d[timeString]; !ok {
			d[timeString] = []htmlMessage{}
		}

		d[timeString] = append(d[timeString], htmlMessage{
			Message:     message,
			Attachments: data.prepareAttachments(conversation.ID, message),
		})
	}

	return d
}

func (data *html) prepareAttachments(conversationID string, message models.SignalMessage) []htmlAttachment {
	d := []htmlAttachment{}
	os.MkdirAll(data.getHTMLAttachmentPath(conversationID), os.ModePerm)

	for _, attachment := range message.Attachments {
		filename := strconv.Itoa(int(attachment.UploadTimestamp))
		if attachment.CdnKey != "" {
			filename = fmt.Sprintf("%s-%s", attachment.CdnKey, filename)
		}
		if attachment.Filename != "" {
			filename = fmt.Sprintf("%s-%s", attachment.Filename, filename)
		}
		filename = strings.ReplaceAll(filename, ".", "-")

		ext, err := attachment.GetExtension()
		if err != nil {
			log.Print(color.RedString("[FAILED] "), fmt.Sprintf("Cannot create extension for attachment %s", attachment.ContentType))
			continue
		}

		var thumbnailPath string
		if attachment.Thumbnail.ContentType != "" {
			thumbnailExt, err := attachment.Thumbnail.GetExtension()
			if err != nil {
				log.Print(color.RedString("[FAILED] "), fmt.Sprintf("Cannot create extension for thumbnail of attachment %s", attachment.Thumbnail.ContentType))
				continue
			}

			thumbnailPath = fmt.Sprintf("%s/%s-thumnail.%s", data.getHTMLAttachmentPath((conversationID)), filename, thumbnailExt)
			err = copyFile(fmt.Sprintf("%s/attachments.noindex/%s", data.signalPath, attachment.Thumbnail.Path), thumbnailPath)

			if err != nil {
				log.Print(color.RedString("[FAILED] "), fmt.Sprintf("Cannot copy thumbnail of attachment %s: %s", attachment.Thumbnail.ContentType, err))
				continue
			}
		}

		path := fmt.Sprintf("%s/%s.%s", data.getHTMLAttachmentPath((conversationID)), filename, ext)
		err = copyFile(fmt.Sprintf("%s/attachments.noindex/%s", data.signalPath, attachment.Path), path)
		if err != nil {
			log.Print(color.RedString("[FAILED] "), fmt.Sprintf("Cannot copy attachment %s: %s", attachment.Thumbnail.ContentType, err))
			continue
		}

		d = append(d, htmlAttachment{
			Attachment:        attachment,
			HTMLPath:          path,
			HTMLThumbnailPath: thumbnailPath,
		})
	}

	return d
}

func copyFile(fromPath, toPath string) error {
	from, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	return err
}
