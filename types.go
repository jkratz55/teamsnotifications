package teams

import "encoding/json"

const (
	// AdaptiveCardContentType is the attachment content type for Adaptive Cards.
	AdaptiveCardContentType = "application/vnd.microsoft.card.adaptive"

	// DefaultAdaptiveCardVersion is used when AdaptiveCard.Version is empty.
	DefaultAdaptiveCardVersion = "1.4"

	// AdaptiveCardSchema is used when AdaptiveCard.Schema is empty.
	AdaptiveCardSchema = "http://adaptivecards.io/schemas/adaptive-card.json"

	// AdaptiveCardWidthFull makes a card span the full width of the message
	// when set as AdaptiveCardMSTeams.Width.
	AdaptiveCardWidthFull = "Full"
)

// The types in this file set their "type" JSON property when marshaled, so it
// never has to be filled in by hand. Each MarshalJSON converts the value to a
// local type without methods to avoid recursing into itself.

// Message is the envelope, with cards as attachments, that Microsoft documents
// for Teams incoming webhooks. Use it when the webhook's flow reads cards from
// the attachments array. Flows that post the request body directly as a card
// expect a bare AdaptiveCard instead. Use NewAdaptiveCardMessage to build one.
type Message struct {
	Attachments []Attachment `json:"attachments"`
}

// NewAdaptiveCardMessage wraps cards in a Message for webhooks that read cards
// from an attachments array.
func NewAdaptiveCardMessage(cards ...AdaptiveCard) Message {
	attachments := make([]Attachment, len(cards))
	for i, card := range cards {
		attachments[i] = Attachment{ContentType: AdaptiveCardContentType, Content: card}
	}
	return Message{Attachments: attachments}
}

// MarshalJSON implements json.Marshaler.
func (m Message) MarshalJSON() ([]byte, error) {
	type message Message
	if m.Attachments == nil {
		m.Attachments = []Attachment{}
	}
	return json.Marshal(struct {
		Type string `json:"type"`
		message
	}{"message", message(m)})
}

// Attachment is a card attached to a Message.
type Attachment struct {
	ContentType string `json:"contentType"`
	Content     any    `json:"content"`
}

// AdaptiveCard represents a top-level Adaptive Card. Post it directly or wrap
// it with NewAdaptiveCardMessage, depending on what the webhook expects.
type AdaptiveCard struct {
	// Version defaults to DefaultAdaptiveCardVersion when empty.
	Version string                `json:"version"`
	Body    []AdaptiveCardElement `json:"body"`
	Actions []AdaptiveCardAction  `json:"actions,omitempty"`
	MSTeams *AdaptiveCardMSTeams  `json:"msteams,omitempty"`
	// Schema defaults to AdaptiveCardSchema when empty.
	Schema string `json:"$schema"`
}

// MarshalJSON implements json.Marshaler.
func (c AdaptiveCard) MarshalJSON() ([]byte, error) {
	type card AdaptiveCard
	if c.Version == "" {
		c.Version = DefaultAdaptiveCardVersion
	}
	if c.Schema == "" {
		c.Schema = AdaptiveCardSchema
	}
	if c.Body == nil {
		c.Body = []AdaptiveCardElement{}
	}
	return json.Marshal(struct {
		Type string `json:"type"`
		card
	}{"AdaptiveCard", card(c)})
}

// AdaptiveCardMSTeams holds Teams-specific Adaptive Card properties.
type AdaptiveCardMSTeams struct {
	// Width set to AdaptiveCardWidthFull makes the card use the full message
	// width instead of the default narrower layout.
	Width string `json:"width,omitempty"`
}

// AdaptiveCardElement is an element that can be placed in a card body,
// container or column. To use an element this package doesn't model, such as
// Image, use RawAdaptiveCardElement.
type AdaptiveCardElement interface {
	adaptiveCardElement()
}

// RawAdaptiveCardElement is the JSON encoding of an element this package doesn't
// model. It is included in the card as is.
type RawAdaptiveCardElement json.RawMessage

func (RawAdaptiveCardElement) adaptiveCardElement() {}

// MarshalJSON implements json.Marshaler.
func (r RawAdaptiveCardElement) MarshalJSON() ([]byte, error) {
	return json.RawMessage(r).MarshalJSON()
}

// TextWeight controls the weight of text in an AdaptiveCardTextBlock.
type TextWeight string

// Supported TextWeight values.
const (
	WeightDefault TextWeight = "Default"
	WeightLighter TextWeight = "Lighter"
	WeightBolder  TextWeight = "Bolder"
)

// TextSize controls the size of text in an AdaptiveCardTextBlock.
type TextSize string

// Supported TextSize values.
const (
	SizeDefault    TextSize = "Default"
	SizeSmall      TextSize = "Small"
	SizeMedium     TextSize = "Medium"
	SizeLarge      TextSize = "Large"
	SizeExtraLarge TextSize = "ExtraLarge"
)

// TextColor controls the color of text in an AdaptiveCardTextBlock.
type TextColor string

// Supported TextColor values.
const (
	ColorDefault   TextColor = "Default"
	ColorDark      TextColor = "Dark"
	ColorLight     TextColor = "Light"
	ColorAccent    TextColor = "Accent"
	ColorGood      TextColor = "Good"
	ColorWarning   TextColor = "Warning"
	ColorAttention TextColor = "Attention"
)

// AdaptiveCardTextBlock is a text item displayed by Adaptive Card renderers.
type AdaptiveCardTextBlock struct {
	Text   string     `json:"text"`
	Weight TextWeight `json:"weight,omitempty"`
	Size   TextSize   `json:"size,omitempty"`
	Color  TextColor  `json:"color,omitempty"`
	Wrap   bool       `json:"wrap,omitempty"`
}

func (AdaptiveCardTextBlock) adaptiveCardElement() {}

// MarshalJSON implements json.Marshaler.
func (t AdaptiveCardTextBlock) MarshalJSON() ([]byte, error) {
	type textBlock AdaptiveCardTextBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		textBlock
	}{"TextBlock", textBlock(t)})
}

// AdaptiveCardFactSet represents a set of title/value facts in an Adaptive Card.
type AdaptiveCardFactSet struct {
	Facts []AdaptiveCardFact `json:"facts"`
}

func (AdaptiveCardFactSet) adaptiveCardElement() {}

// MarshalJSON implements json.Marshaler.
func (f AdaptiveCardFactSet) MarshalJSON() ([]byte, error) {
	type factSet AdaptiveCardFactSet
	if f.Facts == nil {
		f.Facts = []AdaptiveCardFact{}
	}
	return json.Marshal(struct {
		Type string `json:"type"`
		factSet
	}{"FactSet", factSet(f)})
}

// AdaptiveCardFact represents one fact row in an Adaptive Card fact set.
type AdaptiveCardFact struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

// AdaptiveCardContainer groups elements together.
type AdaptiveCardContainer struct {
	Items []AdaptiveCardElement `json:"items"`
}

func (AdaptiveCardContainer) adaptiveCardElement() {}

// MarshalJSON implements json.Marshaler.
func (c AdaptiveCardContainer) MarshalJSON() ([]byte, error) {
	type container AdaptiveCardContainer
	if c.Items == nil {
		c.Items = []AdaptiveCardElement{}
	}
	return json.Marshal(struct {
		Type string `json:"type"`
		container
	}{"Container", container(c)})
}

// AdaptiveCardColumnSet lays out columns side by side.
type AdaptiveCardColumnSet struct {
	Columns []AdaptiveCardColumn `json:"columns"`
}

func (AdaptiveCardColumnSet) adaptiveCardElement() {}

// MarshalJSON implements json.Marshaler.
func (c AdaptiveCardColumnSet) MarshalJSON() ([]byte, error) {
	type columnSet AdaptiveCardColumnSet
	if c.Columns == nil {
		c.Columns = []AdaptiveCardColumn{}
	}
	return json.Marshal(struct {
		Type string `json:"type"`
		columnSet
	}{"ColumnSet", columnSet(c)})
}

// AdaptiveCardColumn is one column of an AdaptiveCardColumnSet.
type AdaptiveCardColumn struct {
	Items []AdaptiveCardElement `json:"items"`
	// Width is "auto", "stretch" or a pixel width such as "50px".
	Width string `json:"width,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (c AdaptiveCardColumn) MarshalJSON() ([]byte, error) {
	type column AdaptiveCardColumn
	if c.Items == nil {
		c.Items = []AdaptiveCardElement{}
	}
	return json.Marshal(struct {
		Type string `json:"type"`
		column
	}{"Column", column(c)})
}

// AdaptiveCardAction is an action shown at the bottom of a card. Teams webhooks
// support Action.OpenUrl, Action.ShowCard and Action.ToggleVisibility. Use
// RawAdaptiveCardAction for actions this package doesn't model.
type AdaptiveCardAction interface {
	adaptiveCardAction()
}

// RawAdaptiveCardAction is the JSON encoding of an action this package doesn't
// model. It is included in the card as is.
type RawAdaptiveCardAction json.RawMessage

func (RawAdaptiveCardAction) adaptiveCardAction() {}

// MarshalJSON implements json.Marshaler.
func (r RawAdaptiveCardAction) MarshalJSON() ([]byte, error) {
	return json.RawMessage(r).MarshalJSON()
}

// AdaptiveCardOpenURLAction is a button that opens a URL.
type AdaptiveCardOpenURLAction struct {
	Title string `json:"title,omitempty"`
	URL   string `json:"url"`
}

func (AdaptiveCardOpenURLAction) adaptiveCardAction() {}

// MarshalJSON implements json.Marshaler.
func (a AdaptiveCardOpenURLAction) MarshalJSON() ([]byte, error) {
	type openURLAction AdaptiveCardOpenURLAction
	return json.Marshal(struct {
		Type string `json:"type"`
		openURLAction
	}{"Action.OpenUrl", openURLAction(a)})
}
