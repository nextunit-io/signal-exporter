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

type SignalMessageExpirationTimer struct {
	ExpireTimer     int    `json:"expireTimer"`
	Source          string `json:"source"`
	FromSync        bool   `json:"fromSync"`
	FromGroupUpdate bool   `json:"fromGroupUpdate"`
}
type SignalMessage struct {
	ID                         string                       `json:"id"`
	Attachments                []SignalAttachments          `json:"attachments"`
	Timestamp                  int                          `json:"timestamp"`
	Source                     string                       `json:"source"`
	SourceUUID                 string                       `json:"sourceUuid"`
	SourceDevice               int                          `json:"sourceDevice"`
	SentAt                     int                          `json:"sent_at"`
	SentTo                     []string                     `json:"sent_to"`
	ServerTimestamp            int                          `json:"serverTimestamp"`
	ReceivedAt                 int                          `json:"received_at"`
	ConversationID             string                       `json:"conversationId"`
	UnidentifiedDeliveryRecord bool                         `json:"unidentifiedDeliveryRecord"`
	Type                       string                       `json:"type"`
	SchemaVersion              int                          `json:"schemaVersion"`
	Body                       string                       `json:"body"`
	DecryptedAt                int                          `json:"decrypted_at"`
	Errors                     []string                     `json:"errors"`
	Flags                      int                          `json:"flags"`
	HasAttachments             int                          `json:"hasAttachments"`
	IsViewOnce                 bool                         `json:"isViewOnce"`
	Preview                    []SignalPreview              `json:"preview"`
	RequiredProtocolVersion    int                          `json:"requiredProtocolVersion"`
	SupportedVersionAtReceive  int                          `json:"supportedVersionAtReceive"`
	Quote                      SignalQuote                  `json:"quote"`
	Sticker                    string                       `json:"sticker"`
	ExpirationTimerUpdate      SignalMessageExpirationTimer `json:"expirationTimerUpdate"`
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
	Timestamp                  string                         `json:"timestamü"`
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
	DraftAttachments           []SignalAttachments            `json:"draftAttachments"`
	DraftTimestamp             int                            `json:"draftTimestamp"`
	isPinned                   bool                           `json:"isPinned"`
	UnreadCount                int                            `json:"unreadCount"`
	Verified                   int                            `json:"verified"`
	MessageCount               int                            `json:"messageCount"`
	SentMessageCount           int                            `json:"sentMessageCount"`
	Messages                   []SignalMessage
}

type SignalData struct {
	Conversations []SignalConverstation
}
