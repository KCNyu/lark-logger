package larklogger

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Note: Constants moved to constants.go for better organization

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
	Actions         []Action `json:"actions,omitempty"`
	FlexMode        string   `json:"flex_mode,omitempty"`
	BackgroundStyle string   `json:"background_style,omitempty"`
	Padding         *Padding `json:"padding,omitempty"`
	TextAlign       string   `json:"text_align,omitempty"`
}

// Action represents button action
type Action struct {
	Tag     string       `json:"tag"`
	Text    *Text        `json:"text,omitempty"`
	URL     string       `json:"url,omitempty"`
	Type    string       `json:"type,omitempty"`
	Value   *ActionValue `json:"value,omitempty"`
	Confirm *Confirm     `json:"confirm,omitempty"`
}

// ActionValue represents action value
type ActionValue struct {
	Key string `json:"key,omitempty"`
}

// Confirm represents confirmation dialog
type Confirm struct {
	Title *Text `json:"title,omitempty"`
	Text  *Text `json:"text,omitempty"`
}

// Button represents a button configuration
type Button struct {
	Text    string `json:"text"`
	URL     string `json:"url,omitempty"`
	Style   string `json:"style,omitempty"`
	Confirm bool   `json:"confirm,omitempty"`
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

// Note: Constants moved to constants.go for better organization
// Background styles
const (
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

// getKeyEmoji returns minimal emoji for key type (simplified for cleaner display)
func getKeyEmoji(key string) string {
	keyLower := strings.ToLower(key)

	switch {
	case strings.Contains(keyLower, "error") || strings.Contains(keyLower, "exception"):
		return EmojiError
	case strings.Contains(keyLower, "warning") || strings.Contains(keyLower, "warn"):
		return EmojiWarn
	case strings.Contains(keyLower, "success") || strings.Contains(keyLower, "ok"):
		return "✅"
	default:
		return "" // No emoji for most fields to reduce visual clutter
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

	// Format long strings for better mobile display
	return formatLongString(valueStr)
}

// formatLongString formats long strings for better mobile display
func formatLongString(value string) string {
	if value == "" {
		return "-"
	}

	// First escape necessary characters
	escaped := escapeMarkdown(value)

	// Handle JSON strings specially
	if strings.HasPrefix(escaped, "{") && strings.HasSuffix(escaped, "}") ||
		strings.HasPrefix(escaped, "[") && strings.HasSuffix(escaped, "]") {
		return formatJSONString(escaped)
	}

	// Handle very long strings
	if len(escaped) > 100 {
		return formatLongTextString(escaped)
	}

	return escaped
}

// formatJSONString formats JSON for better readability
func formatJSONString(jsonStr string) string {
	// For mobile, truncate very long JSON and add length hint, no code fences
	if len(jsonStr) > 200 {
		truncated := jsonStr[:150]
		if lastComma := strings.LastIndex(truncated, ","); lastComma > 50 {
			truncated = jsonStr[:lastComma]
		}
		return truncated + "\n... (" + fmt.Sprintf("%d", len(jsonStr)) + " chars)"
	}
	// For shorter JSON, return as-is (plain text)
	return jsonStr
}

// formatLongTextString formats long text strings
func formatLongTextString(text string) string {
	// Break long strings into multiple lines for mobile readability
	if len(text) > 100 {
		// Insert line breaks at word boundaries
		words := strings.Fields(text)
		var lines []string
		var currentLine string

		for _, word := range words {
			if len(currentLine)+len(word)+1 > 50 { // 50 chars per line for mobile
				if currentLine != "" {
					lines = append(lines, currentLine)
					currentLine = word
				} else {
					// Word itself is too long, break it
					for len(word) > 50 {
						lines = append(lines, word[:50])
						word = word[50:]
					}
					currentLine = word
				}
			} else {
				if currentLine == "" {
					currentLine = word
				} else {
					currentLine += " " + word
				}
			}
		}

		if currentLine != "" {
			lines = append(lines, currentLine)
		}

		return strings.Join(lines, "\n")
	}

	return text
}

// toNonBreaking converts common break chars to non-breaking ones for no-wrap display
func toNonBreaking(s string) string {
	if s == "" {
		return s
	}
	r := strings.ReplaceAll(s, " ", "\u00a0") // space -> NBSP
	r = strings.ReplaceAll(r, "-", "\u2011")  // hyphen -> non-breaking hyphen
	r = strings.ReplaceAll(r, "/", "\u2215")  // slash -> division slash
	return r
}

// shouldStackKV decides if a KV item should switch to vertical block layout
func shouldStackKV(value string) bool {
	if value == "" {
		return false
	}
	if len(value) > 80 {
		return true
	}
	if strings.HasPrefix(value, "{") || strings.HasPrefix(value, "[") {
		return true
	}
	if strings.Contains(value, "\n") {
		return true
	}
	return false
}

// escapeMarkdown escapes only necessary characters, preserving markdown formatting
func escapeMarkdown(content string) string {
	if content == "" {
		return "-"
	}

	// Only escape characters that could break the card structure
	// Preserve markdown formatting like ** for bold and ` for code
	replacer := strings.NewReplacer(
		"<", "&lt;",
		">", "&gt;",
		"&", "&amp;",
	)
	return replacer.Replace(content)
}

// CardBuilder helps build Lark cards
type CardBuilder struct {
	card     *Card
	isMobile bool // Flag for mobile optimization
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
		isMobile: false, // Default to desktop
	}
}

// NewMobileCardBuilder creates a mobile-optimized card builder
func NewMobileCardBuilder() *CardBuilder {
	cb := NewCardBuilder()
	cb.isMobile = true
	return cb
}

// SetMobileOptimization enables/disables mobile optimization
func (cb *CardBuilder) SetMobileOptimization(mobile bool) *CardBuilder {
	cb.isMobile = mobile
	return cb
}

// getPadding returns appropriate padding based on device type
func (cb *CardBuilder) getPadding() *Padding {
	if cb.isMobile {
		return &Padding{
			Top:    MobilePaddingTop,
			Bottom: MobilePaddingBottom,
			Left:   MobilePaddingLeft,
			Right:  MobilePaddingRight,
		}
	}
	return &Padding{
		Top:    PaddingTop,
		Bottom: PaddingBottom,
		Left:   PaddingLeft,
		Right:  PaddingRight,
	}
}

// getLineHeight returns appropriate line height based on device type
func (cb *CardBuilder) getLineHeight() string {
	if cb.isMobile {
		return MobileLineHeight
	}
	return LineHeight
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
	// Use mobile-optimized padding if mobile flag is set
	padding := &Padding{
		Top:    PaddingTop,
		Bottom: PaddingBottom,
		Left:   PaddingLeft,
		Right:  PaddingRight,
	}

	lineHeight := LineHeight
	if cb.isMobile {
		padding = &Padding{
			Top:    MobilePaddingTop,
			Bottom: MobilePaddingBottom,
			Left:   MobilePaddingLeft,
			Right:  MobilePaddingRight,
		}
		lineHeight = MobileLineHeight
	}

	cb.card.Card.Elements = append(cb.card.Card.Elements, Element{
		Tag: "div",
		Text: &Text{
			Tag:        "lark_md",
			Content:    escapeMarkdown(content),
			LineHeight: lineHeight,
		},
		Padding:   padding,
		TextAlign: "left",
	})
	return cb
}

// AddSubtitle adds subtitle with message
func (cb *CardBuilder) AddSubtitle(subtitle string) *CardBuilder {
	// Use mobile-optimized padding if mobile flag is set
	padding := &Padding{
		Top:    PaddingTop,
		Bottom: 0,
		Left:   PaddingLeft,
		Right:  PaddingRight,
	}

	lineHeight := LineHeight
	if cb.isMobile {
		padding = &Padding{
			Top:    MobilePaddingTop,
			Bottom: 0,
			Left:   MobilePaddingLeft,
			Right:  MobilePaddingRight,
		}
		lineHeight = MobileLineHeight
	}

	cb.card.Card.Elements = append(cb.card.Card.Elements, Element{
		Tag: "div",
		Text: &Text{
			Tag:        "lark_md",
			Content:    fmt.Sprintf("<font color=\"grey\">%s</font>", escapeMarkdown(subtitle)),
			LineHeight: lineHeight,
		},
		Padding:   padding,
		TextAlign: "left",
	})
	return cb
}

// AddTimestamp adds timestamp (right-aligned)
func (cb *CardBuilder) AddTimestamp() *CardBuilder {
	// Use mobile-optimized padding if mobile flag is set
	padding := &Padding{
		Top:    0,
		Bottom: PaddingBottom,
		Left:   PaddingLeft,
		Right:  PaddingRight,
	}

	lineHeight := LineHeight
	textAlign := "right"

	if cb.isMobile {
		padding = &Padding{
			Top:    0,
			Bottom: MobilePaddingBottom,
			Left:   MobilePaddingLeft,
			Right:  MobilePaddingRight,
		}
		lineHeight = MobileLineHeight
		textAlign = "left" // Left align on mobile for better readability
	}

	cb.card.Card.Elements = append(cb.card.Card.Elements, Element{
		Tag: "div",
		Text: &Text{
			Tag:        "lark_md",
			Content:    fmt.Sprintf("<font color=\"grey\">%s %s</font>", EmojiTime, time.Now().Format("2006-01-02 15:04:05")),
			LineHeight: lineHeight,
		},
		Padding:   padding,
		TextAlign: textAlign,
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
	// Add section title with emoji, bold formatting and center alignment
	cb.card.Card.Elements = append(cb.card.Card.Elements, Element{
		Tag: "div",
		Text: &Text{
			Tag:        "lark_md",
			Content:    "📊 **Data Fields**",
			LineHeight: cb.getLineHeight(),
		},
		Padding:   cb.getPadding(),
		TextAlign: "center", // Center align for better visual hierarchy
	})

	// Split into short and long groups
	var shortItems []KVItem
	var longItems []KVItem
	for _, kv := range kvList {
		if shouldStackKV(kv.Value) {
			longItems = append(longItems, kv)
		} else {
			shortItems = append(shortItems, kv)
		}
	}

	// Render short items as table (two-column rows)
	if len(shortItems) > 0 {
		// Add table header with clean styling (no emojis), edge-to-edge
		headerColumns := []Column{
			{
				Tag:           "column",
				Width:         "weighted",
				Weight:        ColumnWeightKey,
				VerticalAlign: "middle",
				Elements: []ColumnElement{
					{Tag: "markdown", Content: "**Key**", TextAlign: "left", FontSize: FontSizeLarge},
				},
			},
			{
				Tag:           "column",
				Width:         "weighted",
				Weight:        ColumnWeightValue,
				VerticalAlign: "middle",
				Elements: []ColumnElement{
					{Tag: "markdown", Content: "**Value**", TextAlign: "left", FontSize: FontSizeLarge},
				},
			},
		}
		cb.card.Card.Elements = append(cb.card.Card.Elements, Element{
			Tag:             "column_set",
			Columns:         headerColumns,
			FlexMode:        "none",
			BackgroundStyle: "grey",
			Padding:         &Padding{Top: PaddingTop, Bottom: PaddingBottom, Left: 0, Right: 0},
		})

		// Rows
		for i, kv := range shortItems {
			bgStyle := "default"
			if i%2 == 1 {
				bgStyle = "light"
			}
			keyNoWrap := toNonBreaking(kv.Key)
			valueDisplay := kv.Value
			row := []Column{
				{Tag: "column", Width: "weighted", Weight: ColumnWeightKey, VerticalAlign: "top", Elements: []ColumnElement{{Tag: "markdown", Content: "**" + keyNoWrap + "**", TextAlign: "left", FontSize: FontSizeDefault}}},
				{Tag: "column", Width: "weighted", Weight: ColumnWeightValue, VerticalAlign: "top", Elements: []ColumnElement{{Tag: "markdown", Content: valueDisplay, TextAlign: "left", FontSize: FontSizeDefault}}},
			}
			cb.card.Card.Elements = append(cb.card.Card.Elements, Element{Tag: "column_set", Columns: row, FlexMode: "none", BackgroundStyle: bgStyle, Padding: &Padding{Top: PaddingTop / 2, Bottom: PaddingBottom / 2, Left: 0, Right: 0}})
		}
	}

	// Render long items as a separate block
	if len(longItems) > 0 {
		cb.AddDivider()
		for _, kv := range longItems {
			keyNoWrap := toNonBreaking(kv.Key)
			cb.card.Card.Elements = append(cb.card.Card.Elements, Element{
				Tag:       "div",
				Text:      &Text{Tag: "lark_md", Content: "**" + keyNoWrap + "**\n" + kv.Value, LineHeight: cb.getLineHeight()},
				Padding:   &Padding{Top: PaddingTop / 2, Bottom: PaddingBottom / 2, Left: 0, Right: 0},
				TextAlign: "left",
			})
		}
	}

	return cb
}

// AddConfigGrid adds a compact 2x2 configuration grid
func (cb *CardBuilder) AddConfigGrid(configData map[string]string) *CardBuilder {
	// Cross layout: two columns; labels row then values row for each pair group
	compact := &Padding{Top: 2, Bottom: 2, Left: 0, Right: 0}

	// Row 1: labels (Level, Service)
	rowLabels1 := []Column{
		{ // Level label
			Tag:           "column",
			Width:         "weighted",
			Weight:        1,
			VerticalAlign: "top",
			Elements: []ColumnElement{{
				Tag:       "markdown",
				Content:   fmt.Sprintf("**%s %s:**", "📊", toNonBreaking("Level")),
				TextAlign: "left",
				FontSize:  FontSizeSmall,
			}},
		},
		{ // Service label
			Tag:           "column",
			Width:         "weighted",
			Weight:        1,
			VerticalAlign: "top",
			Elements: []ColumnElement{{
				Tag:       "markdown",
				Content:   fmt.Sprintf("**%s %s:**", "🔧", toNonBreaking("Service")),
				TextAlign: "left",
				FontSize:  FontSizeSmall,
			}},
		},
	}
	cb.card.Card.Elements = append(cb.card.Card.Elements, Element{Tag: "column_set", Columns: rowLabels1, FlexMode: "none", BackgroundStyle: ColorLightBlue, Padding: compact})

	// Row 2: values (Level value, Service value)
	rowValues1 := []Column{
		{
			Tag:           "column",
			Width:         "weighted",
			Weight:        1,
			VerticalAlign: "top",
			Elements: []ColumnElement{{
				Tag:       "markdown",
				Content:   configData["level_value"],
				TextAlign: "left",
				FontSize:  FontSizeDefault,
			}},
		},
		{
			Tag:           "column",
			Width:         "weighted",
			Weight:        1,
			VerticalAlign: "top",
			Elements: []ColumnElement{{
				Tag:       "markdown",
				Content:   configData["service_value"],
				TextAlign: "left",
				FontSize:  FontSizeDefault,
			}},
		},
	}
	cb.card.Card.Elements = append(cb.card.Card.Elements, Element{Tag: "column_set", Columns: rowValues1, FlexMode: "none", BackgroundStyle: ColorLightBlue, Padding: compact})

	// Row 3: labels (Env, Hostname)
	rowLabels2 := []Column{
		{ // Env label
			Tag:           "column",
			Width:         "weighted",
			Weight:        1,
			VerticalAlign: "top",
			Elements: []ColumnElement{{
				Tag:       "markdown",
				Content:   fmt.Sprintf("**%s %s:**", "🌍", toNonBreaking("Env")),
				TextAlign: "left",
				FontSize:  FontSizeSmall,
			}},
		},
		{ // Hostname label
			Tag:           "column",
			Width:         "weighted",
			Weight:        1,
			VerticalAlign: "top",
			Elements: []ColumnElement{{
				Tag:       "markdown",
				Content:   fmt.Sprintf("**%s %s:**", "🖥️", toNonBreaking("Hostname")),
				TextAlign: "left",
				FontSize:  FontSizeSmall,
			}},
		},
	}
	cb.card.Card.Elements = append(cb.card.Card.Elements, Element{Tag: "column_set", Columns: rowLabels2, FlexMode: "none", BackgroundStyle: ColorLightBlue, Padding: compact})

	// Row 4: values (Env value, Hostname value)
	rowValues2 := []Column{
		{
			Tag:           "column",
			Width:         "weighted",
			Weight:        1,
			VerticalAlign: "top",
			Elements: []ColumnElement{{
				Tag:       "markdown",
				Content:   configData["env_value"],
				TextAlign: "left",
				FontSize:  FontSizeDefault,
			}},
		},
		{
			Tag:           "column",
			Width:         "weighted",
			Weight:        1,
			VerticalAlign: "top",
			Elements: []ColumnElement{{
				Tag:       "markdown",
				Content:   configData["hostname_value"],
				TextAlign: "left",
				FontSize:  FontSizeDefault,
			}},
		},
	}
	cb.card.Card.Elements = append(cb.card.Card.Elements, Element{Tag: "column_set", Columns: rowValues2, FlexMode: "none", BackgroundStyle: ColorLightBlue, Padding: compact})

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

// AddButtons adds action buttons to the card
func (cb *CardBuilder) AddButtons(buttons []Button) *CardBuilder {
	if len(buttons) == 0 {
		return cb
	}

	var actions []Action
	for _, button := range buttons {
		action := Action{
			Tag:  "button",
			Text: &Text{Tag: "plain_text", Content: button.Text},
			Type: "url",
			URL:  button.URL,
		}

		// Set button style with enhanced colors for confirm actions
		if button.Style != "" {
			action.Type = button.Style
		}

		// Use danger style for confirm buttons to indicate caution
		if button.Confirm {
			action.Type = ButtonStyleDanger
			action.Confirm = &Confirm{
				Title: &Text{Tag: "plain_text", Content: "⚠️ Confirm Action"},
				Text:  &Text{Tag: "plain_text", Content: fmt.Sprintf("Are you sure you want to execute %s?\n\nThis action cannot be undone.", button.Text)},
			}
		} else if button.Style == ButtonStylePrimary {
			action.Type = ButtonStylePrimary
		} else {
			action.Type = ButtonStyleSecondary
		}

		actions = append(actions, action)
	}

	cb.card.Card.Elements = append(cb.card.Card.Elements, Element{
		Tag:     "action",
		Actions: actions,
	})

	return cb
}

// AddButton adds a single button to the card
func (cb *CardBuilder) AddButton(text, url string, style ...string) *CardBuilder {
	buttonStyle := ButtonStyleSecondary
	if len(style) > 0 {
		buttonStyle = style[0]
	}

	button := Button{
		Text:  text,
		URL:   url,
		Style: buttonStyle,
	}

	return cb.AddButtons([]Button{button})
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
