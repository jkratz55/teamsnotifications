package teams

// AdaptiveCard represents a top-level Adaptive Card payload that can be sent to
// webhook endpoints that expect Adaptive Card JSON.
type AdaptiveCard struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Body    []any  `json:"body"`
	Schema  string `json:"$schema"`
}

// AdaptiveCardTextBlock is a text item displayed by Adaptive Card renderers.
type AdaptiveCardTextBlock struct {
	Type   string `json:"type"`
	Text   string `json:"text"`
	Weight string `json:"weight,omitempty"`
	Size   string `json:"size,omitempty"`
	Color  string `json:"color,omitempty"`
	Wrap   bool   `json:"wrap,omitempty"`
}

// AdaptiveCardFactSet represents a set of title/value facts in an Adaptive Card.
type AdaptiveCardFactSet struct {
	Type  string             `json:"type"`
	Facts []AdaptiveCardFact `json:"facts"`
}

// AdaptiveCardFact represents one fact row in an Adaptive Card fact set.
type AdaptiveCardFact struct {
	Title string `json:"title"`
	Value string `json:"value"`
}
