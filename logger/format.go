package logger

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func itoa(n int) string { return strconv.Itoa(n) }

// appendText renders e as a single human readable line:
//
//	15:04:05.000  INFO   bot.polling  message  key=value
func appendText(b []byte, e Entry, timeFmt string, color bool) []byte {
	paint := func(style, s string) {
		if color && style != "" {
			b = append(b, style...)
			b = append(b, s...)
			b = append(b, ansiReset...)
		} else {
			b = append(b, s...)
		}
	}
	if timeFmt != "" {
		paint(ansiGray, e.Time.Format(timeFmt))
		b = append(b, ' ', ' ')
	}
	badge, textStyle := levelStyle(e.Level)
	if color {
		paint(badge, " "+e.Level.short()+" ")
	} else {
		b = append(b, e.Level.short()...)
	}
	b = append(b, ' ', ' ')
	if e.Name != "" {
		paint(ansiBold+ansiMagenta, e.Name)
		b = append(b, ' ', ' ')
	}
	msg := e.Message
	if strings.ContainsRune(msg, '\n') {
		msg = strings.ReplaceAll(msg, "\n", "\n    ")
	}
	paint(textStyle, msg)
	for _, f := range e.Fields {
		if f.Key == "" || f.Value == nil {
			continue
		}
		if err, ok := f.Value.(error); ok && err == nil {
			continue
		}
		b = append(b, ' ', ' ')
		paint(ansiCyan, f.Key)
		paint(ansiGray, "=")
		v := fieldString(f.Value)
		if _, isErr := f.Value.(error); isErr {
			paint(ansiBrightRed, v)
		} else {
			paint("", v)
		}
	}
	if e.Caller != "" {
		b = append(b, ' ', ' ')
		paint(ansiGray, "("+e.Caller+")")
	}
	return append(b, '\n')
}

// fieldString formats a field value, quoting it when it would be
// ambiguous on a key=value line.
func fieldString(v any) string {
	var s string
	switch t := v.(type) {
	case string:
		s = t
	case error:
		s = t.Error()
	case fmt.Stringer:
		s = t.String()
	default:
		s = fmt.Sprint(v)
	}
	if s == "" || strings.IndexFunc(s, func(r rune) bool {
		return unicode.IsSpace(r) || r == '=' || r == '"'
	}) >= 0 {
		return strconv.Quote(s)
	}
	return s
}

// appendJSON renders e as a JSON object followed by a newline.
func appendJSON(b []byte, e Entry) []byte {
	b = append(b, `{"time":`...)
	b = strconv.AppendQuote(b[:len(b)], e.Time.UTC().Format("2006-01-02T15:04:05.000Z07:00"))
	b = append(b, `,"level":`...)
	b = strconv.AppendQuote(b, strings.ToLower(e.Level.String()))
	if e.Name != "" {
		b = append(b, `,"logger":`...)
		b = appendJSONString(b, e.Name)
	}
	b = append(b, `,"msg":`...)
	b = appendJSONString(b, e.Message)
	if e.Caller != "" {
		b = append(b, `,"caller":`...)
		b = appendJSONString(b, e.Caller)
	}
	for _, f := range e.Fields {
		if f.Key == "" || f.Value == nil {
			continue
		}
		b = append(b, ',')
		b = appendJSONString(b, f.Key)
		b = append(b, ':')
		switch t := f.Value.(type) {
		case error:
			b = appendJSONString(b, t.Error())
		case fmt.Stringer:
			b = appendJSONString(b, t.String())
		default:
			enc, err := json.Marshal(f.Value)
			if err != nil {
				enc, _ = json.Marshal(fmt.Sprint(f.Value))
			}
			b = append(b, enc...)
		}
	}
	return append(b, '}', '\n')
}

func appendJSONString(b []byte, s string) []byte {
	enc, _ := json.Marshal(s)
	return append(b, enc...)
}
