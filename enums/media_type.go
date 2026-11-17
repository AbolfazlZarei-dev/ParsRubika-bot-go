package enums

// MediaType is the kind of inline attachment a message carries.
type MediaType string

const (
	MediaTypeText       MediaType = "Text"
	MediaTypeFile       MediaType = "File"
	MediaTypeFileInline MediaType = "FileInline"
	MediaTypeSticker    MediaType = "Sticker"
	MediaTypeCaption    MediaType = "FileInlineCaption"
	MediaTypeEvent      MediaType = "Event"
)

// String implements fmt.Stringer.
func (m MediaType) String() string { return string(m) }
