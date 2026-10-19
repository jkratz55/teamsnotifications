package teams

import (
	"encoding/json"
	"testing"
)

func TestMarshalSetsTypeAndDefaults(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{
			name:  "empty message",
			value: Message{},
			want:  `{"type":"message","attachments":[]}`,
		},
		{
			name:  "empty card",
			value: AdaptiveCard{},
			want:  `{"type":"AdaptiveCard","version":"1.4","body":[],"$schema":"http://adaptivecards.io/schemas/adaptive-card.json"}`,
		},
		{
			name:  "card keeps explicit version and schema",
			value: AdaptiveCard{Version: "1.5", Schema: "https://example.com/schema.json"},
			want:  `{"type":"AdaptiveCard","version":"1.5","body":[],"$schema":"https://example.com/schema.json"}`,
		},
		{
			name:  "pointer element",
			value: AdaptiveCard{Body: []AdaptiveCardElement{&AdaptiveCardTextBlock{Text: "hi"}}},
			want:  `{"type":"AdaptiveCard","version":"1.4","body":[{"type":"TextBlock","text":"hi"}],"$schema":"http://adaptivecards.io/schemas/adaptive-card.json"}`,
		},
		{
			name:  "empty fact set",
			value: AdaptiveCardFactSet{},
			want:  `{"type":"FactSet","facts":[]}`,
		},
		{
			name: "container and column set",
			value: AdaptiveCardContainer{Items: []AdaptiveCardElement{
				AdaptiveCardColumnSet{Columns: []AdaptiveCardColumn{
					{Width: "auto", Items: []AdaptiveCardElement{AdaptiveCardTextBlock{Text: "a"}}},
					{},
				}},
			}},
			want: `{"type":"Container","items":[{"type":"ColumnSet","columns":[{"type":"Column","items":[{"type":"TextBlock","text":"a"}],"width":"auto"},{"type":"Column","items":[]}]}]}`,
		},
		{
			name:  "open URL action",
			value: AdaptiveCardOpenURLAction{Title: "Open", URL: "https://example.com"},
			want:  `{"type":"Action.OpenUrl","title":"Open","url":"https://example.com"}`,
		},
		{
			name: "raw element and action",
			value: AdaptiveCard{
				Body:    []AdaptiveCardElement{RawAdaptiveCardElement(`{"type":"Image","url":"https://example.com/a.png"}`)},
				Actions: []AdaptiveCardAction{RawAdaptiveCardAction(`{"type":"Action.ToggleVisibility","targetElements":["a"]}`)},
			},
			want: `{"type":"AdaptiveCard","version":"1.4","body":[{"type":"Image","url":"https://example.com/a.png"}],"actions":[{"type":"Action.ToggleVisibility","targetElements":["a"]}],"$schema":"http://adaptivecards.io/schemas/adaptive-card.json"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.value)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("json.Marshal() =\n%s\nwant\n%s", got, tt.want)
			}
		})
	}
}

func TestMarshalRejectsInvalidRawElement(t *testing.T) {
	card := AdaptiveCard{Body: []AdaptiveCardElement{RawAdaptiveCardElement(`{not json`)}}
	if _, err := json.Marshal(card); err == nil {
		t.Error("json.Marshal() error = nil, want error for invalid raw element")
	}
}

func TestNewAdaptiveCardMessage(t *testing.T) {
	msg := NewAdaptiveCardMessage(AdaptiveCard{}, AdaptiveCard{Version: "1.5"})
	if len(msg.Attachments) != 2 {
		t.Fatalf("len(Attachments) = %d, want 2", len(msg.Attachments))
	}
	for i, a := range msg.Attachments {
		if a.ContentType != AdaptiveCardContentType {
			t.Errorf("Attachments[%d].ContentType = %q, want %q", i, a.ContentType, AdaptiveCardContentType)
		}
	}
}
