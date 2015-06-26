package robots

type SlashCommand struct {
	Payload
	Command string `schema:"command"`
}

type Payload struct {
	Token       string  `schema:"token"`
	TeamID      string  `schema:"team_id"`
	TeamDomain  string  `schema:"team_domain,omitempty"`
	ChannelID   string  `schema:"channel_id"`
	ChannelName string  `schema:"channel_name"`
	Timestamp   float64 `schema:"timestamp,omitempty"`
	UserID      string  `schema:"user_id"`
	UserName    string  `schema:"user_name"`
	Text        string  `schema:"text,omitempty"`
	TriggerWord string  `schema:"trigger_word,omitempty"`
	Service_ID  string  `schema:"service_id,omitempty"`
	Robot       string
}

type OutgoingWebHook struct {
	Payload
	TriggerWord string `schema:"trigger_word"`
}

type OutgoingWebHookResponse struct {
	Text      string     `json:"text"`
	Parse     ParseStyle `json:"parse,omitempty"`
	LinkNames bool       `json:"link_names,omitempty"`
	Markdown  bool       `json:"mrkdwn,omitempty"`
}

type ParseStyle string

var (
	ParseStyleFull = ParseStyle("full")
	ParseStyleNone = ParseStyle("none")
)

type IncomingWebhook struct {
	Domain      string       `json:"domain"`
	Channel     string       `json:"channel"`
	Username    string       `json:"username"`
	Text        string       `json:"text"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	UnfurlLinks bool         `json:"unfurl_links,omitempty"`
	Parse       ParseStyle   `json:"parse,omitempty"`
	LinkNames   bool         `json:"link_names,omitempty"`
	Markdown    bool         `json:"mrkdwn,omitempty"`
}

type Attachment struct {
	Fallback   string            `json:"fallback,omitempty"`
	Pretext    string            `json:"pretext,omitempty"`
	Text       string            `json:"text,omitempty"`
	Color      string            `json:"color,omitempty"`
	Fields     []AttachmentField `json:"fields,omitempty"`
	MarkdownIn []MarkdownField   `json:"mrkdown_in,omitempty"`
	ImageURL   string            `json:"image_url,omitempty"`
}

type MarkdownField string

var (
	MarkdownFieldPretext  = MarkdownField("pretext")
	MarkdownFieldText     = MarkdownField("text")
	MarkdownFieldTitle    = MarkdownField("title")
	MarkdownFieldFields   = MarkdownField("fields")
	MarkdownFieldFallback = MarkdownField("fallback")
)

type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short,omitempty"`
}

type UserResp struct {
	Ok   bool    `json:"ok"`
	User UserObj `json:"user,omitempty"`
}

type UserObj struct {
	Id                string     `json:"id,omitempty"`
	Name              string     `json:"name,omitempty"`
	Deleted           bool       `json:"deleted,omitempty"`
	Status            string     `json:"status,omitempty,omitempty"`
	Color             string     `json:"color,omitempty"`
	RealName          string     `json:"real_name,omitempty"`
	Tz                string     `json:"tz,omitempty"`
	TzLabel           string     `json:"tz_label,omitempty"`
	TzOffset          int        `json:"tz_offset,omitempty"`
	Profile           ProfileObj `json:"profile,omitempty"`
	IsAdmin           bool       `json:"is_admin,omitempty"`
	IsOwner           bool       `json:"is_owner,omitempty"`
	IsPrimaryOwner    bool       `json:"is_primary_owner,omitempty"`
	IsRestricted      bool       `json:"is_restrcited,omitempty"`
	IsUltraRestricted bool       `json:"is_ultra_restricted,omitempty"`
	IsBot             bool       `json:"is_bot,omitempty"`
	HasFiles          bool       `json:"has_files,omitempty"`
	Has2fa            bool       `json:"has_2fa,omitempty"`
}

type ProfileObj struct {
	FirstName          string `json:"first_name,omitempty"`
	LastName           string `json:"last_name,omitempty"`
	Title              string `json:"title,omitempty"`
	Skype              string `json:"skype,omitempty"`
	Phone              string `json:"phone,omitempty"`
	Image24            string `json:"image_24,omitempty"`
	Image32            string `json:"image_32,omitempty"`
	Image48            string `json:"image_48,omitempty"`
	Image72            string `json:"image_72,omitempty"`
	Image192           string `json:"image_192,omitempty"`
	ImageOriginal      string `json:"image_original,omitempty"`
	RealName           string `json:"real_name,omitempty"`
	RealNameNormalized string `json:"real_name_normalized,omitempty"`
	Email              string `json:"email,omitempty"`
}

type Configuration struct {
	Port   int                   `json:"port"`
	Tokens map[string]Credential `json:"tokens"`
}

type Credential struct {
	IncomingWebhookToken string `json:"incoming_webhook_token,omitempty"`
	APIToken             string `json:"api_token,omitempty"`
}

type Robot interface {
	Run(p *Payload) (botString string)
	Description() (description string)
}
