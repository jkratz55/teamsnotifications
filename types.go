package teams

import (
	"strings"
)

// Message represents the payload structure for creating a notification message in
// Microsoft Teams via webhook.
type Message struct {
	Title    string     `json:"title"`
	Summary  string     `json:"summary,omitempty"`
	Text     string     `json:"text,omitempty"`
	Theme    ThemeColor `json:"themeColor,omitempty"`
	Sections []Section  `json:"sections,omitempty"`
}

// Section represents a logical grouping in a notification message to organize content.
//
// ActivityTitle defines the title displayed in the section header.
// ActivitySubtitle specifies an optional subtitle for the section.
// Facts is a collection of key-value pairs providing additional information in the section.
// Markdown indicates if the section content supports markdown formatting.
type Section struct {
	ActivityTitle    string `json:"activityTitle,omitempty"`
	ActivitySubtitle string `json:"activitySubtitle,omitempty"`
	Facts            []Fact `json:"facts,omitempty"`
	Markdown         bool   `json:"markdown"`
}

// Fact represents a key-value pair, typically used to store related data or
// additional informational content.
type Fact struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ThemeColor represents a string type used to define color codes for various UI
// themes such as success, warning, error, etc.
type ThemeColor string

const (

	// ThemeDefault represents the default theme color, which is no theme.
	ThemeDefault ThemeColor = ""

	// ThemeSuccess represents the color code for the success theme, typically
	// used to indicate successful operations.
	ThemeSuccess ThemeColor = "#19e013"

	// ThemeWarning represents the color code for the warning theme, typically
	// used to indicate caution or alerts.
	ThemeWarning ThemeColor = "#f3b911"

	// ThemeError represents the color code for the error theme, typically used
	// to indicate errors or critical issues.
	ThemeError ThemeColor = "#f40909"

	// ThemeInfo represents the color code for the informational theme, typically
	// used to convey general information.
	ThemeInfo ThemeColor = "#1951fa"
)

// ThemeColorFromHex converts a hex color code string to a ThemeColor type.
// If the input is empty, it returns the default theme color.
// If the input is missing a '#' prefix, it adds one before conversion.
func ThemeColorFromHex(hexColor string) ThemeColor {
	if strings.TrimSpace(hexColor) == "" {
		return ThemeDefault
	}
	if hexColor != "" && hexColor[0] != '#' {
		hexColor = "#" + hexColor
	}
	return ThemeColor(hexColor)
}
