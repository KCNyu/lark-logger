package larklogger

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Emoji constants
const (
	EmojiInfo    = "ℹ️"
	EmojiWarn    = "⚠️"
	EmojiError   = "❌"
	EmojiDefault = "📋"
)

// Color constants
const (
	ColorBlue   = "blue"
	ColorGreen  = "green"
	ColorOrange = "orange"
	ColorRed    = "red"
	ColorGrey   = "grey"
)

// Card represents a Lark interactive card
type Card struct {
	MsgType string   `json:"msg_type"`
	Card    CardData `json:"card"`
}

// CardData represents the card structure
type CardData struct {
	Config       Config    `json:"config"`
	Header       Header    `json:"header"`
	Elements     []Element `json:"elements"`
	CardLink     *CardLink `json:"card_link,omitempty"`
	CornerRadius int       `json:"corner_radius,omitempty"`
}

// Config represents card configuration
type Config struct {
	WideScreenMode bool          `json:"wide_screen_mode"`
	EnableForward  bool          `json:"enable_forward,omitempty"`
	UpdateMulti    bool          `json:"update_multi,omitempty"`
	IosConfig      *MobileConfig `json:"ios_config,omitempty"`
	AndroidConfig  *MobileConfig `json:"android_config,omitempty"`
}

// MobileConfig represents mobile configuration
type MobileConfig struct {
	EnableForward bool `json:"enable_forward,omitempty"`
}

// Header represents card header
type Header struct {
	Title    Title  `json:"title"`
	SubTitle *Title `json:"subtitle,omitempty"`
	Template string `json:"template"`
}

// Title represents title structure
type Title struct {
	Tag      string `json:"tag"`
	Content  string `json:"content"`
	FontSize string `json:"font_size,omitempty"`
}

// Element represents card element
type Element struct {
	Tag             string   `json:"tag"`
	Text            *Text    `json:"text,omitempty"`
	Columns         []Column `json:"columns,omitempty"`
	FlexMode        string   `json:"flex_mode,omitempty"`
	BackgroundStyle string   `json:"background_style,omitempty"`
	Padding         *Padding `json:"padding,omitempty"`
	TextAlign       string   `json:"text_align,omitempty"`
}

// Padding represents padding structure
type Padding struct {
	Top    int `json:"top,omitempty"`
	Bottom int `json:"bottom,omitempty"`
	Left   int `json:"left,omitempty"`
	Right  int `json:"right,omitempty"`
}

// Column represents column structure
type Column struct {
	Tag           string          `json:"tag"`                      // Fixed "column"
	Width         string          `json:"width"`                    // Fixed "weighted"
	Weight        int             `json:"weight,omitempty"`         // Key=3, Value=7 (3:7 layout)
	VerticalAlign string          `json:"vertical_align,omitempty"` // Fixed "middle"
	Elements      []ColumnElement `json:"elements,omitempty"`       // Column content
}

// ColumnElement represents column element
type ColumnElement struct {
	Tag       string `json:"tag"`                  // Fixed "markdown"
	Content   string `json:"content"`              // Text content
	TextAlign string `json:"text_align,omitempty"` // Text alignment
	FontSize  string `json:"font_size,omitempty"`  // Font size for important values
}

// Text represents text element
type Text struct {
	Tag        string `json:"tag"`                   // lark_md
	Content    string `json:"content"`               // Text content
	LineHeight string `json:"line_height,omitempty"` // Line height: 1.5
}

// CardLink represents card link
type CardLink struct {
	URL string `json:"url"` // Target URL
}

// CardField represents a field in the card
type CardField struct {
	IsShort bool      `json:"is_short"`
	Text    *CardText `json:"text"`
}

// CardText represents text in card elements
type CardText struct {
	Tag     string `json:"tag"`     // plain_text or lark_md
	Content string `json:"content"` // Text content
}

// CardSection represents a section in the card
type CardSection struct {
	Tag    string      `json:"tag"`
	Text   *CardText   `json:"text,omitempty"`
	Fields []*CardText `json:"fields,omitempty"`
}

// -------------------------- Visual Constants --------------------------

const (
	// Corner radius (px)
	CornerRadius = 8
	// Padding (px): top/bottom 8, left/right 12
	PaddingTop        = 8
	PaddingBottom     = 8
	PaddingLeft       = 12
	PaddingRight      = 12
	ColumnWeightKey   = 3
	ColumnWeightValue = 7
	// Font sizes
	FontSizeDefault = "default"
	FontSizeLarge   = "large"
	// Line height
	LineHeight = "1.5"
	// Background styles
	BgStyleHeader = "grey"    // Header: light grey
	BgStyleOdd    = "default" // Odd rows: white
	BgStyleEven   = "light"   // Even rows: very light grey
)

// getVisualConfig returns visual configuration based on log level
func getVisualConfig(level LogLevel) (template string) {
	switch level {
	case LevelInfo:
		return ColorBlue
	case LevelWarn:
		return ColorOrange
	case LevelError:
		return ColorRed
	default:
		return ColorGrey
	}
}

// getLogLevelEmoji returns emoji for log level
func getLogLevelEmoji(level LogLevel) string {
	switch level {
	case LevelInfo:
		return EmojiInfo
	case LevelWarn:
		return EmojiWarn
	case LevelError:
		return EmojiError
	default:
		return EmojiDefault
	}
}

// getKeyEmoji returns emoji for key type
func getKeyEmoji(key string) string {
	keyLower := strings.ToLower(key)

	switch {
	case strings.Contains(keyLower, "error") || strings.Contains(keyLower, "exception"):
		return EmojiError
	case strings.Contains(keyLower, "warning") || strings.Contains(keyLower, "warn"):
		return EmojiWarn
	case strings.Contains(keyLower, "success") || strings.Contains(keyLower, "ok"):
		return "✅"
	case strings.Contains(keyLower, "time") || strings.Contains(keyLower, "duration"):
		return "⏱️"
	case strings.Contains(keyLower, "memory") || strings.Contains(keyLower, "ram"):
		return "💾"
	case strings.Contains(keyLower, "cpu") || strings.Contains(keyLower, "processor"):
		return "🖥️"
	case strings.Contains(keyLower, "network") || strings.Contains(keyLower, "connection"):
		return "🌐"
	case strings.Contains(keyLower, "database") || strings.Contains(keyLower, "db"):
		return "🗄️"
	case strings.Contains(keyLower, "user") || strings.Contains(keyLower, "client"):
		return "👤"
	case strings.Contains(keyLower, "request") || strings.Contains(keyLower, "api"):
		return "📡"
	case strings.Contains(keyLower, "file") || strings.Contains(keyLower, "path"):
		return "📁"
	case strings.Contains(keyLower, "port") || strings.Contains(keyLower, "address"):
		return "🔌"
	case strings.Contains(keyLower, "version") || strings.Contains(keyLower, "v"):
		return "🏷️"
	case strings.Contains(keyLower, "count") || strings.Contains(keyLower, "number"):
		return "🔢"
	case strings.Contains(keyLower, "size") || strings.Contains(keyLower, "length"):
		return "📏"
	case strings.Contains(keyLower, "status") || strings.Contains(keyLower, "state"):
		return "📊"
	default:
		return EmojiDefault
	}
}

// KVItem represents a key-value item
type KVItem struct {
	Key   string // Original key
	Value string // Processed value
}

// mapToKVItems converts map to KV items
func mapToKVItems(data map[string]interface{}) []KVItem {
	var items []KVItem

	for k, v := range data {
		if k == "" {
			continue
		}
		items = append(items, KVItem{
			Key:   k,
			Value: formatValue(v),
		})
	}

	return items
}

// formatValue formats value (supports multiple types)
func formatValue(v interface{}) string {
	if v == nil {
		return "-"
	}

	// Handle different value types
	var valueStr string
	switch val := v.(type) {
	case string:
		valueStr = val
	case int, int8, int16, int32, int64:
		valueStr = fmt.Sprintf("%d", val)
	case uint, uint8, uint16, uint32, uint64:
		valueStr = fmt.Sprintf("%d", val)
	case float32, float64:
		valueStr = fmt.Sprintf("%.2f", val)
	case bool:
		valueStr = fmt.Sprintf("%t", val)
	case time.Time:
		valueStr = val.Format("2006-01-02 15:04:05")
	default:
		// Other types (slices, structs): JSON serialization
		jsonBytes, err := json.Marshal(val)
		if err != nil {
			valueStr = fmt.Sprintf("%v", val)
		} else {
			valueStr = string(jsonBytes)
		}
	}

	// Escape special characters to avoid format breaking
	return escapeMarkdown(valueStr)
}

// escapeMarkdown escapes Lark Markdown special characters
func escapeMarkdown(content string) string {
	if content == "" {
		return "-"
	}

	// Simple escape for common markdown characters
	replacer := strings.NewReplacer(
		"*", "\\*",
		"_", "\\_",
		"`", "\\`",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"<", "&lt;",
		">", "&gt;",
		"&", "&amp;",
	)
	return replacer.Replace(content)
}

// CardBuilder helps build Lark cards
type CardBuilder struct {
	card *Card
}

// NewCardBuilder creates a new card builder
func NewCardBuilder() *CardBuilder {
	return &CardBuilder{
		card: &Card{
			MsgType: "interactive",
			Card: CardData{
				Config: Config{
					WideScreenMode: true,
					EnableForward:  true,
					UpdateMulti:    true,
					IosConfig: &MobileConfig{
						EnableForward: true,
					},
					AndroidConfig: &MobileConfig{
						EnableForward: true,
					},
				},
				CornerRadius: CornerRadius,
			},
		},
	}
}

// SetHeader sets the card header
func (cb *CardBuilder) SetHeader(title, template string) *CardBuilder {
	cb.card.Card.Header = Header{
		Title: Title{
			Tag:      "lark_md",
			Content:  escapeMarkdown(title),
			FontSize: FontSizeLarge,
		},
		Template: template,
	}
	return cb
}

// AddSection adds a text section
func (cb *CardBuilder) AddSection(content string) *CardBuilder {
	cb.card.Card.Elements = append(cb.card.Card.Elements, Element{
		Tag: "div",
		Text: &Text{
			Tag:        "lark_md",
			Content:    escapeMarkdown(content),
			LineHeight: LineHeight,
		},
		Padding: &Padding{
			Top:    PaddingTop,
			Bottom: PaddingBottom,
			Left:   PaddingLeft,
			Right:  PaddingRight,
		},
		TextAlign: "left",
	})
	return cb
}

// AddSubtitle adds subtitle with message
func (cb *CardBuilder) AddSubtitle(subtitle string) *CardBuilder {
	cb.card.Card.Elements = append(cb.card.Card.Elements, Element{
		Tag: "div",
		Text: &Text{
			Tag:        "lark_md",
			Content:    fmt.Sprintf("<font color=\"grey\">%s</font>", escapeMarkdown(subtitle)),
			LineHeight: LineHeight,
		},
		Padding: &Padding{
			Top:    PaddingTop,
			Bottom: 0,
			Left:   PaddingLeft,
			Right:  PaddingRight,
		},
		TextAlign: "left",
	})
	return cb
}

// AddTimestamp adds timestamp (right-aligned)
func (cb *CardBuilder) AddTimestamp() *CardBuilder {
	cb.card.Card.Elements = append(cb.card.Card.Elements, Element{
		Tag: "div",
		Text: &Text{
			Tag:        "lark_md",
			Content:    fmt.Sprintf("<font color=\"grey\">⏰ %s</font>", time.Now().Format("2006-01-02 15:04:05")),
			LineHeight: LineHeight,
		},
		Padding: &Padding{
			Top:    0,
			Bottom: PaddingBottom,
			Left:   PaddingLeft,
			Right:  PaddingRight,
		},
		TextAlign: "right",
	})
	return cb
}

// AddDivider adds a divider line
func (cb *CardBuilder) AddDivider() *CardBuilder {
	cb.card.Card.Elements = append(cb.card.Card.Elements, Element{Tag: "hr"})
	return cb
}

// AddKVTable adds professional KV table with alternating colors
func (cb *CardBuilder) AddKVTable(kvList []KVItem) *CardBuilder {
	// Add section title
	cb.AddSection("📊 Data Fields")

	// Add table header with enhanced styling
	headerColumns := []Column{
		{
			Tag:           "column",
			Width:         "weighted",
			Weight:        ColumnWeightKey,
			VerticalAlign: "middle",
			Elements: []ColumnElement{
				{
					Tag:       "markdown",
					Content:   "🔑 **Key**",
					TextAlign: "left",
					FontSize:  FontSizeLarge,
				},
			},
		},
		{
			Tag:           "column",
			Width:         "weighted",
			Weight:        ColumnWeightValue,
			VerticalAlign: "middle",
			Elements: []ColumnElement{
				{
					Tag:       "markdown",
					Content:   "💎 **Value**",
					TextAlign: "left",
					FontSize:  FontSizeLarge,
				},
			},
		},
	}
	cb.card.Card.Elements = append(cb.card.Card.Elements, Element{
		Tag:             "column_set",
		Columns:         headerColumns,
		FlexMode:        "none",
		BackgroundStyle: "grey",
		Padding: &Padding{
			Top:    PaddingTop,
			Bottom: PaddingBottom,
			Left:   PaddingLeft,
			Right:  PaddingRight,
		},
	})

	// Add data rows with enhanced styling
	for i, kv := range kvList {
		// Enhanced alternating background colors
		bgStyle := "default"
		if i%2 == 1 {
			bgStyle = "light"
		}

		dataColumns := []Column{
			{
				Tag:           "column",
				Width:         "weighted",
				Weight:        ColumnWeightKey,
				VerticalAlign: "middle",
				Elements: []ColumnElement{
					{
						Tag:       "markdown",
						Content:   fmt.Sprintf("**%s**", kv.Key),
						TextAlign: "left",
					},
				},
			},
			{
				Tag:           "column",
				Width:         "weighted",
				Weight:        ColumnWeightValue,
				VerticalAlign: "middle",
				Elements: []ColumnElement{
					{
						Tag:       "markdown",
						Content:   kv.Value,
						TextAlign: "left",
						FontSize:  FontSizeDefault,
					},
				},
			},
		}

		cb.card.Card.Elements = append(cb.card.Card.Elements, Element{
			Tag:             "column_set",
			Columns:         dataColumns,
			FlexMode:        "none",
			BackgroundStyle: bgStyle,
			Padding: &Padding{
				Top:    PaddingTop,
				Bottom: PaddingBottom,
				Left:   PaddingLeft,
				Right:  PaddingRight,
			},
		})
	}
	return cb
}

// AddKVTableWithStyle adds a key-value table with custom background style
func (cb *CardBuilder) AddKVTableWithStyle(kvList []KVItem, bgStyle string) *CardBuilder {
	// Add table header
	headerColumns := []Column{
		{
			Tag:           "column",
			Width:         "weighted",
			Weight:        ColumnWeightKey,
			VerticalAlign: "middle",
			Elements: []ColumnElement{
				{
					Tag:       "markdown",
					Content:   "**Key**",
					TextAlign: "left",
					FontSize:  FontSizeLarge,
				},
			},
		},
		{
			Tag:           "column",
			Width:         "weighted",
			Weight:        ColumnWeightValue,
			VerticalAlign: "middle",
			Elements: []ColumnElement{
				{
					Tag:       "markdown",
					Content:   "**Value**",
					TextAlign: "left",
					FontSize:  FontSizeLarge,
				},
			},
		},
	}
	cb.card.Card.Elements = append(cb.card.Card.Elements, Element{
		Tag:             "column_set",
		Columns:         headerColumns,
		FlexMode:        "none",
		BackgroundStyle: bgStyle,
		Padding: &Padding{
			Top:    PaddingTop,
			Bottom: PaddingBottom,
			Left:   PaddingLeft,
			Right:  PaddingRight,
		},
	})

	// Add data rows with custom background style
	for i, kv := range kvList {
		// Use custom background style for all rows
		rowBgStyle := bgStyle
		if i%2 == 1 {
			// Alternate with a slightly different shade
			rowBgStyle = "light_grey"
		}

		rowColumns := []Column{
			{
				Tag:           "column",
				Width:         "weighted",
				Weight:        ColumnWeightKey,
				VerticalAlign: "middle",
				Elements: []ColumnElement{
					{
						Tag:       "markdown",
						Content:   fmt.Sprintf("**%s**", kv.Key),
						TextAlign: "left",
						FontSize:  FontSizeDefault,
					},
				},
			},
			{
				Tag:           "column",
				Width:         "weighted",
				Weight:        ColumnWeightValue,
				VerticalAlign: "middle",
				Elements: []ColumnElement{
					{
						Tag:       "markdown",
						Content:   kv.Value,
						TextAlign: "left",
						FontSize:  FontSizeDefault,
					},
				},
			},
		}
		cb.card.Card.Elements = append(cb.card.Card.Elements, Element{
			Tag:             "column_set",
			Columns:         rowColumns,
			FlexMode:        "none",
			BackgroundStyle: rowBgStyle,
			Padding: &Padding{
				Top:    PaddingTop,
				Bottom: PaddingBottom,
				Left:   PaddingLeft,
				Right:  PaddingRight,
			},
		})
	}
	return cb
}

// AddConfigGrid adds a 2x2 configuration grid with emojis
func (cb *CardBuilder) AddConfigGrid(config map[string]string) *CardBuilder {
	// Create 2x2 grid layout
	gridColumns := []Column{
		{
			Tag:           "column",
			Width:         "weighted",
			Weight:        1,
			VerticalAlign: "middle",
			Elements: []ColumnElement{
				{
					Tag:       "markdown",
					Content:   fmt.Sprintf("**%s**\n%s", config["level"], config["level_value"]),
					TextAlign: "center",
					FontSize:  FontSizeDefault,
				},
			},
		},
		{
			Tag:           "column",
			Width:         "weighted",
			Weight:        1,
			VerticalAlign: "middle",
			Elements: []ColumnElement{
				{
					Tag:       "markdown",
					Content:   fmt.Sprintf("**%s**\n%s", config["service"], config["service_value"]),
					TextAlign: "center",
					FontSize:  FontSizeDefault,
				},
			},
		},
	}

	// First row
	cb.card.Card.Elements = append(cb.card.Card.Elements, Element{
		Tag:             "column_set",
		Columns:         gridColumns,
		FlexMode:        "none",
		BackgroundStyle: "light_blue",
		Padding: &Padding{
			Top:    PaddingTop,
			Bottom: PaddingBottom,
			Left:   PaddingLeft,
			Right:  PaddingRight,
		},
	})

	// Second row
	gridColumns2 := []Column{
		{
			Tag:           "column",
			Width:         "weighted",
			Weight:        1,
			VerticalAlign: "middle",
			Elements: []ColumnElement{
				{
					Tag:       "markdown",
					Content:   fmt.Sprintf("**%s**\n%s", config["env"], config["env_value"]),
					TextAlign: "center",
					FontSize:  FontSizeDefault,
				},
			},
		},
		{
			Tag:           "column",
			Width:         "weighted",
			Weight:        1,
			VerticalAlign: "middle",
			Elements: []ColumnElement{
				{
					Tag:       "markdown",
					Content:   fmt.Sprintf("**%s**\n%s", config["hostname"], config["hostname_value"]),
					TextAlign: "center",
					FontSize:  FontSizeDefault,
				},
			},
		},
	}

	cb.card.Card.Elements = append(cb.card.Card.Elements, Element{
		Tag:             "column_set",
		Columns:         gridColumns2,
		FlexMode:        "none",
		BackgroundStyle: "light_blue",
		Padding: &Padding{
			Top:    PaddingTop,
			Bottom: PaddingBottom,
			Left:   PaddingLeft,
			Right:  PaddingRight,
		},
	})

	return cb
}

// AddKeyValueList adds a key-value list section
func (cb *CardBuilder) AddKeyValueList(title string, kv map[string]interface{}) *CardBuilder {
	cb.AddSection(fmt.Sprintf("**%s**", title))

	for key, value := range kv {
		formattedValue := formatValue(value)
		cb.AddSection(fmt.Sprintf("**%s**: %s", key, formattedValue))
	}

	return cb
}

// AddStatusBadge adds a status badge
func (cb *CardBuilder) AddStatusBadge(status, message string) *CardBuilder {
	emoji := "✅"
	if status == "error" || status == "failed" {
		emoji = EmojiError
	} else if status == "warning" || status == "warn" {
		emoji = EmojiWarn
	}

	cb.AddSection(fmt.Sprintf("%s **Status**: %s", emoji, message))
	return cb
}

// AddMetricsGrid adds a metrics grid
func (cb *CardBuilder) AddMetricsGrid(title string, metrics map[string]interface{}) *CardBuilder {
	cb.AddSection(fmt.Sprintf("**%s**", title))

	// Create simple list for metrics
	var contents []string
	for key, value := range metrics {
		formattedValue := formatValue(value)
		contents = append(contents, fmt.Sprintf("**%s**: %s", key, formattedValue))
	}

	cb.card.Card.Elements = append(cb.card.Card.Elements, Element{
		Tag: "div",
		Text: &Text{
			Tag:        "lark_md",
			Content:    strings.Join(contents, "\n"),
			LineHeight: LineHeight,
		},
		Padding: &Padding{
			Top:    PaddingTop,
			Bottom: PaddingBottom,
			Left:   PaddingLeft,
			Right:  PaddingRight,
		},
		TextAlign: "left",
	})

	return cb
}

// AddCardLink adds optional card link
func (cb *CardBuilder) AddCardLink(url string) *CardBuilder {
	if url != "" {
		cb.card.Card.CardLink = &CardLink{URL: url}
		cb.card.Card.Elements = append(cb.card.Card.Elements, Element{
			Tag: "div",
			Text: &Text{
				Tag:        "lark_md",
				Content:    fmt.Sprintf("<font color=\"grey\">📌 Click card to view detailed logs: [Log Link](%s)</font>", escapeMarkdown(url)),
				LineHeight: LineHeight,
			},
			Padding: &Padding{
				Top:    PaddingTop,
				Bottom: PaddingBottom,
				Left:   PaddingLeft,
				Right:  PaddingRight,
			},
			TextAlign: "left",
		})
	}
	return cb
}

// Build builds the card
func (cb *CardBuilder) Build() *Card {
	return cb.card
}

// ToJSON converts to JSON
func (c *Card) ToJSON() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to marshal card to JSON: %w", err)
	}
	return string(data), nil
}

// NewCardField creates a new card field
func NewCardField(isShort bool, content string) *CardField {
	return &CardField{
		IsShort: isShort,
		Text: &CardText{
			Tag:     "lark_md",
			Content: content,
		},
	}
}

// FormatTimestamp formats timestamp for display
func FormatTimestamp(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// GetLogLevelEmoji returns emoji for log level
func GetLogLevelEmoji(level LogLevel) string {
	switch level {
	case LevelInfo:
		return EmojiInfo
	case LevelWarn:
		return EmojiWarn
	case LevelError:
		return EmojiError
	default:
		return EmojiDefault
	}
}
