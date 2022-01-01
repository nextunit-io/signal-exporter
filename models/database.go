package models

import (
	"fmt"
)

// TODO Bodyranges
type SignalAttachment struct {
	*SignalImage

	Screenshot SignalMediaBase `json:"screenshot"`
}

type SignalPreview struct {
	Url         string      `json:"url"`
	Title       string      `json:"title"`
	Image       SignalImage `json:"image"`
	Description string      `json:"description"`
	Date        interface{} `json:"date"`
}

type SignalMediaBase struct {
	ContentType string `json:"contentType"`
	Path        string `json:"path"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
}

type SignalMediaType int

const (
	UNKNOWN SignalMediaType = iota
	VIDEO
	IMAGE
	AUDIO
	TEXT
)

type SignalImage struct {
	*SignalMediaBase

	AttachmentIdentifier string          `json:"attachment_identifier"`
	CdnKey               string          `json:"cdnKey"`
	Size                 int             `json:"size"`
	Filename             string          `json:"fileName"`
	Flags                int             `json:"flags"`
	Caption              string          `json:"caption"`
	BlurHash             string          `json:"blurHash"`
	UploadTimestamp      int64           `json:"uploadTimestamp`
	CdnNumber            int             `json:"cdnNumber"`
	Thumbnail            SignalMediaBase `json:"thumbnail"`
}

type SignalQuote struct {
	ID          interface{}        `json:"id"`
	Author      string             `json:"author"`
	AuthorUUID  string             `json:"authorUuid"`
	Text        string             `json:"text"`
	Attachments []SignalAttachment `json:"attachments"`
}

type SignalMessageExpirationTimer struct {
	ExpireTimer     int    `json:"expireTimer"`
	Source          string `json:"source"`
	FromSync        bool   `json:"fromSync"`
	FromGroupUpdate bool   `json:"fromGroupUpdate"`
}
type SignalMessage struct {
	ID                         string                       `json:"id"`
	Attachments                []SignalAttachment           `json:"attachments"`
	Timestamp                  int64                        `json:"timestamp"`
	Source                     string                       `json:"source"`
	SourceUUID                 string                       `json:"sourceUuid"`
	SourceDevice               int                          `json:"sourceDevice"`
	SentAt                     int                          `json:"sent_at"`
	SentTo                     []string                     `json:"sent_to"`
	ServerTimestamp            int64                        `json:"serverTimestamp"`
	ReceivedAt                 int                          `json:"received_at"`
	ConversationID             string                       `json:"conversationId"`
	UnidentifiedDeliveryRecord bool                         `json:"unidentifiedDeliveryRecord"`
	Type                       string                       `json:"type"`
	SchemaVersion              int                          `json:"schemaVersion"`
	Body                       string                       `json:"body"`
	DecryptedAt                int                          `json:"decrypted_at"`
	Errors                     []interface{}                `json:"errors"`
	Flags                      int                          `json:"flags"`
	HasAttachments             int                          `json:"hasAttachments"`
	IsViewOnce                 bool                         `json:"isViewOnce"`
	Preview                    []SignalPreview              `json:"preview"`
	RequiredProtocolVersion    int                          `json:"requiredProtocolVersion"`
	SupportedVersionAtReceive  int                          `json:"supportedVersionAtReceive"`
	Quote                      SignalQuote                  `json:"quote"`
	Sticker                    SignalSticker                `json:"sticker"`
	ExpirationTimerUpdate      SignalMessageExpirationTimer `json:"expirationTimerUpdate"`
}

type SignalSticker struct {
	PackId    string            `json:"packId"`
	PackKey   string            `json:"packKey"`
	StickerId int               `json:"stickerId"`
	Data      SignalStickerData `json:"data"`
}

type SignalStickerData struct {
	ID          int         `json:"id"`
	PackId      string      `json:"packId"`
	Emoji       string      `json:"emoji"`
	IsCoverOnly bool        `json:"isCoverOnly"`
	LastUsed    interface{} `json:"lastUsed"`
	Path        string      `json:"path"`
	Height      int         `json:"height"`
	Width       int         `json:"width"`
	ContentType string      `json:"contentType"`
	Size        int         `json:"size"`
}

type SignalProfileAvatar struct {
	Hash string `json:"hash"`
	Path string `json:"path"`
}

type SignalProfile struct {
	ProfileAvatar        SignalProfileAvatar `json:"profileAvatar"`
	ProfileKey           string              `json:"profileKey"`
	ProfileKeyVersion    string              `json:"profileKeyVersion"`
	ProfileKeyCredential string              `json:"profileKeyCredential"`
	AccessKey            string              `json:"accessKey"`
	ProfileName          string              `json:"profileName"`
	StorageProfileKey    string              `json:"storageProfileKey"`
	StorageID            string              `json:"storageID"`
	ProfileSharing       bool                `json:"profileSharing"`
	Name                 string              `json:"name"`
}

type SignalConversationCapabilities struct {
	Gv2          bool `json:"gv2"`
	Gv1Migration bool `json:"gv1-migration"`
}
type SignalConverstation struct {
	*SignalProfile

	ID                         string                         `json:"id"`
	UUID                       string                         `json:"uuid"`
	E164                       string                         `json:"e164"`
	GroupID                    string                         `json:"groupId"`
	Type                       string                         `json:"type"`
	Version                    int                            `json:"version"`
	SealedSender               int                            `json:"sealedSender"`
	ActiveAt                   int                            `json:"active_at"`
	Archived                   bool                           `json:"isArchived"`
	LastMessage                string                         `json:"lastMessage"`
	LastMessageStatus          string                         `json:"lastMessageStatus"`
	Timestamp                  int64                          `json:"timestamü"`
	Capabilities               SignalConversationCapabilities `json:"capabilities"`
	MarkedUnread               bool                           `json:"markedUnread"`
	SharedGroupNames           []string                       `json:"sharedGroupNames"`
	ExpireTimer                int                            `json:"expireTimer"`
	MessageRequestResponseType int                            `json:"messageRequestResponseType"`
	Draft                      string                         `json:"draft"`
	DraftChanged               bool                           `json:"draftChanged"`
	DraftBodyRanges            []string                       `json:"draftBodyRanges"`
	Color                      string                         `json:"color"`
	InboxPosition              int                            `json:"inbox_position"`
	Avatar                     SignalProfileAvatar            `json:"avatar"`
	QuoteMessageID             string                         `json:"quoteMessageId"`
	DraftAttachments           []SignalAttachment             `json:"draftAttachments"`
	DraftTimestamp             int64                          `json:"draftTimestamp"`
	IsPinned                   bool                           `json:"isPinned"`
	UnreadCount                int                            `json:"unreadCount"`
	Verified                   int                            `json:"verified"`
	MessageCount               int                            `json:"messageCount"`
	SentMessageCount           int                            `json:"sentMessageCount"`
	Messages                   []SignalMessage
}

type SignalData struct {
	Conversations []SignalConverstation
}

func (file *SignalMediaBase) GetFileType() (SignalMediaType, error) {
	switch file.ContentType {
	case "image/gif":
		return IMAGE, nil
	case "image/jpeg":
		return IMAGE, nil
	case "image/png":
		return IMAGE, nil
	case "image/tiff":
		return IMAGE, nil
	case "image/svg+xml":
		return IMAGE, nil
	case "audio/mpeg":
		return AUDIO, nil
	case "audio/aac":
		return AUDIO, nil
	case "text/css":
		return TEXT, nil
	case "text/csv":
		return TEXT, nil
	case "text/html":
		return TEXT, nil
	case "text/javascript":
		return TEXT, nil
	case "text/plain":
		return TEXT, nil
	case "text/xml":
		return TEXT, nil
	case "application/pdf":
		return TEXT, nil
	case "video/mpeg":
		return VIDEO, nil
	case "video/mp4":
		return VIDEO, nil
	case "video/quicktime":
		return VIDEO, nil
	}

	return UNKNOWN, fmt.Errorf("Not supported type %s", file.ContentType)
}

func (file *SignalMediaBase) GetExtension() (string, error) {
	switch file.ContentType {
	case "image/gif":
		return "gif", nil
	case "image/jpeg":
		return "jpeg", nil
	case "image/png":
		return "png", nil
	case "image/tiff":
		return "tiff", nil
	case "image/svg+xml":
		return "svg", nil
	case "audio/mpeg":
		return "mpeg", nil
	case "audio/aac":
		return "aac", nil
	case "text/css":
		return "css", nil
	case "text/csv":
		return "csv", nil
	case "text/html":
		return "html", nil
	case "text/javascript":
		return "js", nil
	case "text/plain":
		return "txt", nil
	case "text/xml":
		return "xml", nil
	case "video/mpeg":
		return "mpeg", nil
	case "video/mp4":
		return "mp4", nil
	case "video/quicktime":
		return "mp4", nil
	case "application/pdf":
		return "pdf", nil
	}

	return "", fmt.Errorf("Not supported type %s", file.ContentType)
}
