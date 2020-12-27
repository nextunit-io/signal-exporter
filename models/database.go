package models

// TODO Bodyranges
// TODO Attachments
type SignalAttachments struct {
}

type SignalPreview struct {
	Url         string      `json:"url"`
	Title       string      `json:"title"`
	Image       SignalImage `json:"image"`
	Description string      `json:"description"`
	Date        string      `json:"date"`
}

type SignalImageBase struct {
	ContentType string `json:"contentType"`
	Path        string `json:"path"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
}

type SignalImage struct {
	*SignalImageBase

	AttachmentIdentifier string          `json:"attachment_identifier"`
	CdnKey               string          `json:"cdnKey"`
	Size                 int             `json:"size"`
	Filename             string          `json:"fileName"`
	Flags                string          `json:"flags"`
	Caption              string          `json:"caption"`
	BlurHash             string          `json:"blurHash"`
	UploadTimestamp      int             `json:"uploadTimestamp`
	CdnNumber            int             `json:"cdnNumber"`
	Thumbnail            SignalImageBase `json:"thumbnail"`
}

type SignalQuote struct {
	ID          int                 `json:"id"`
	Author      string              `json:"author"`
	AuthorUUID  string              `json:"authorUuid"`
	Text        string              `json:"text"`
	Attachments []SignalAttachments `json:"attachments"`
}

type SignalMessage struct {
	ID                         string              `json:"id"`
	Attachments                []SignalAttachments `json:"attachments"`
	Timestamp                  int                 `json:"timestamp"`
	Source                     string              `json:"source"`
	SourceUUID                 string              `json:"sourceUuid"`
	SourceDevice               int                 `json:"sourceDevice"`
	SentAt                     int                 `json:"sent_at"`
	SentTo                     []string            `json:"sent_to"`
	ServerTimestamp            int                 `json:"serverTimestamp"`
	ReceivedAt                 int                 `json:"received_at"`
	ConversationID             string              `json:"conversationId"`
	UnidentifiedDeliveryRecord bool                `json:"unidentifiedDeliveryRecord"`
	Type                       string              `json:"type"`
	SchemaVersion              int                 `json:"schemaVersion"`
	Body                       string              `json:"body"`
	DecryptedAt                int                 `json:"decrypted_at"`
	Errors                     []string            `json:"errors"`
	Flags                      int                 `json:"flags"`
	HasAttachments             int                 `json:"hasAttachments"`
	IsViewOnce                 bool                `json:"isViewOnce"`
	Preview                    []SignalPreview     `json:"preview"`
	RequiredProtocolVersion    int                 `json:"requiredProtocolVersion"`
	SupportedVersionAtReceive  int                 `json:"supportedVersionAtReceive"`
	Quote                      SignalQuote         `json:"quote"`
	Sticker                    string              `json:"sticker"`
}

type SignalConverstation struct {
	ConversationID string
	Messages       []SignalMessage
}

type SignalData struct {
	Conversations map[string]*SignalConverstation
}
