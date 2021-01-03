package html

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"nextunit.io/signal-exporter/models"
)

const indexTemplatePath = "html/index.html"
const messageTemplatePath = "html/message.html"

var attachmentLinkFunc = func(path string) string {
	return fmt.Sprintf("./attachments/%s", filepath.Base(path))
}

var templateFunctions = template.FuncMap{
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

		return attachmentLinkFunc(path)
	},
	"attachmentLink": attachmentLinkFunc,
	"attachmentIsVideo": func(attachment htmlAttachment) bool {
		attachmentType, err := attachment.Attachment.GetFileType()
		if err != nil {
			log.Print(color.RedString("[FAILED] "), fmt.Sprintf("attachmentIsVideo is failing for type %s: %s", attachment.Attachment.ContentType, err))
			return false
		}

		return attachmentType == models.VIDEO
	},
}

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
	templateFunctions["generateMessagePath"] = func(conversationID, date string) string {
		return strings.ReplaceAll(fmt.Sprintf("%s/%s.html", data.getHTMLConversationPath(conversationID), date), fmt.Sprintf("%s/", data.getHTMLPath()), "")
	}

	err := data.generateMessages()
	if err != nil {
		return err
	}

	return data.generateIndex()
}

func (data *html) generateIndex() error {
	files, err := ioutil.ReadDir(data.getHTMLPath())
	if err != nil {
		log.Print(color.RedString("[FAILED] "), fmt.Sprintf("cannot generate index - html dir not readable: %s", err))
		return err
	}

	conversations := []*models.SignalConverstation{}
	dates := map[string][]string{}

	for _, file := range files {
		if file.IsDir() {
			conversationID := file.Name()
			dates[conversationID] = []string{}

			jsonFile, err := os.Open(fmt.Sprintf("%s/%s", data.getHTMLConversationPath(conversationID), data.getHTMLConversationConfigFilename()))
			if err != nil {
				log.Print(color.RedString("[FAILED] "), fmt.Sprintf("cannot generate index - json file: %s", err))
				return err
			}

			byteValue, _ := ioutil.ReadAll(jsonFile)
			err = jsonFile.Close()
			if err != nil {
				log.Print(color.RedString("[FAILED] "), fmt.Sprintf("cannot close file: %s", err))
				return err
			}

			var conversation models.SignalConverstation
			err = json.Unmarshal(byteValue, &conversation)
			if err != nil {
				log.Print(color.RedString("[FAILED] "), fmt.Sprintf("cannot unmarshal: %s", err))
				return err
			}

			if conversation.ProfileAvatar.Path != "" {
				err = copyFile(fmt.Sprintf("%s/attachments.noindex/%s", data.signalPath, conversation.ProfileAvatar.Path), fmt.Sprintf("%s/profile.png", data.getHTMLConversationPath(conversationID)))
				if err != nil {
					log.Print(color.RedString("[FAILED] "), fmt.Sprintf("cannot generate profile avatar: %s", err))
					return err
				}
			}

			conversations = append(conversations, &conversation)

			messageFiles, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", data.getHTMLPath(), file.Name()))
			if err != nil {
				log.Print(color.RedString("[FAILED] "), fmt.Sprintf("Cannot generate index - message files not readable: %s", err))
				return err
			}

			for _, messageFile := range messageFiles {
				if !messageFile.IsDir() && path.Ext(messageFile.Name()) == ".html" {
					matched, err := regexp.Match(`^\d{4}-\d{2}-\d{2}`, []byte(messageFile.Name()))

					if matched && err == nil {
						dates[conversationID] = append(dates[conversationID], strings.TrimSuffix(messageFile.Name(), path.Ext(messageFile.Name())))
					}
				}
			}
		}
	}

	t := template.Must(template.New(filepath.Base(indexTemplatePath)).Funcs(templateFunctions).ParseFiles(indexTemplatePath))
	f, err := os.Create(fmt.Sprintf("%s/index.html", data.getHTMLPath()))
	if err != nil {
		log.Print(color.RedString("[FAILED] "), fmt.Sprintf("Cannot create index file: %s", err))
		return err
	}

	err = t.Execute(f, struct {
		Conversations []*models.SignalConverstation
		Dates         map[string][]string
	}{
		Conversations: conversations,
		Dates:         dates,
	})

	if err != nil {
		log.Print(color.RedString("[FAILED] "), fmt.Sprintf("Cannot execute index generation: %s", err))
		return err
	}
	return nil
}

func (data *html) generateMessages() error {
	counter := 0
	for _, c := range data.data.Conversations {
		if c.MessageCount == 0 {
			continue
		}
		counter++

		t := template.Must(template.New(filepath.Base(messageTemplatePath)).Funcs(templateFunctions).ParseFiles(messageTemplatePath))
		os.MkdirAll(data.getHTMLConversationPath(c.ID), os.ModePerm)

		conversation := c
		conversation.Messages = []models.SignalMessage{}

		jsonData, err := json.Marshal(conversation)
		if err != nil {
			log.Print(color.RedString("[FAILED] "), fmt.Sprintf("Cannot generate config file: %s", err))
			return err
		}
		err = ioutil.WriteFile(data.getHTMLConversationConfigFilePath(c.ID), jsonData, 0644)
		if err != nil {
			log.Print(color.RedString("[FAILED] "), fmt.Sprintf("Cannot generate config file: %s", err))
			return err
		}

		messages := data.prepareMessages(c)

		for dateString, message := range messages {
			sort.Sort(SortMessages(message))

			f, err := os.Create(fmt.Sprintf("%s/%s.html", data.getHTMLConversationPath(c.ID), dateString))
			if err != nil {
				log.Println("create file: ", err)
				return err
			}

			err = t.Execute(f, struct {
				Conversation models.SignalConverstation
				Messages     []htmlMessage
				Config       struct {
					MediaTypeVideo models.SignalMediaType
				}
			}{
				Conversation: c,
				Messages:     message,
				Config: struct {
					MediaTypeVideo models.SignalMediaType
				}{
					MediaTypeVideo: models.VIDEO,
				},
			})

			if err != nil {
				log.Print(color.RedString("[FAILED] "), fmt.Sprintf("Cannot execute template generation: %s", err))
				return err
			}
		}
	}

	log.Print(color.GreenString("[OK] "), fmt.Sprintf("Generated %d directories", counter))
	return nil
}

func (data *html) getHTMLConversationPath(conversationID string) string {
	return fmt.Sprintf("%s/%s", data.getHTMLPath(), conversationID)
}
func (data *html) getHTMLPath() string {
	return fmt.Sprintf("%s/html", data.outputPath)
}
func (data *html) getHTMLAttachmentPath(conversationID string) string {
	return fmt.Sprintf("%s/attachments", data.getHTMLConversationPath(conversationID))
}
func (data *html) getHTMLConversationConfigFilePath(conversationID string) string {
	return fmt.Sprintf("%s/%s", data.getHTMLConversationPath(conversationID), data.getHTMLConversationConfigFilename())
}
func (data *html) getHTMLConversationConfigFilename() string {
	return "config.json"
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
