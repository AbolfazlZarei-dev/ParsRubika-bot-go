package widgets

import (
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

func TestChatKeypadBuilder(t *testing.T) {
	k := NewChat().
		AddButton(Simple("btn_1", "One")).
		Row().
		AddButton(Simple("btn_2", "Two")).
		AddButton(Simple("btn_3", "Three")).
		Build()
	if k == nil {
		t.Fatal("Build returned nil")
	}
	if len(k.Rows) != 2 {
		t.Errorf("rows = %d, want 2", len(k.Rows))
	}
	if len(k.Rows[0].Buttons) != 1 {
		t.Errorf("row 0 buttons = %d, want 1", len(k.Rows[0].Buttons))
	}
	if len(k.Rows[1].Buttons) != 2 {
		t.Errorf("row 1 buttons = %d, want 2", len(k.Rows[1].Buttons))
	}
	if k.Rows[1].Buttons[0].ID != "btn_2" {
		t.Errorf("row 1 button 0 ID = %q", k.Rows[1].Buttons[0].ID)
	}
}

func TestInlineBuilder(t *testing.T) {
	k := NewInline().
		AddButton(Simple("btn_1", "One")).
		AddButton(Simple("btn_2", "Two")).
		Row().
		AddButton(Link("btn_link", "Visit", "https://example.com")).
		Build()
	if k == nil {
		t.Fatal("Build returned nil")
	}
	if len(k.Rows) != 2 {
		t.Errorf("rows = %d, want 2", len(k.Rows))
	}
	if k.Rows[1].Buttons[0].ButtonLink == nil {
		t.Error("Link button should have ButtonLink set")
	}
}

func TestJoinChannelButton(t *testing.T) {
	b := JoinChannel("btn", "Join", "@mychannel")
	if b.ButtonLink == nil {
		t.Fatal("ButtonLink should not be nil")
	}
	if b.ButtonLink.JoinChannelData == nil {
		t.Fatal("JoinChannelData should not be nil")
	}
	if b.ButtonLink.JoinChannelData.Username != "mychannel" {
		t.Errorf("username = %q", b.ButtonLink.JoinChannelData.Username)
	}
}

func TestOpenChatButton(t *testing.T) {
	b := OpenChat("btn", "Open", "g0GUID", "Group")
	if b.ButtonLink == nil || b.ButtonLink.OpenChatData == nil {
		t.Fatal("OpenChatData should not be nil")
	}
	if b.ButtonLink.OpenChatData.ObjectGUID != "g0GUID" {
		t.Errorf("ObjectGUID = %q", b.ButtonLink.OpenChatData.ObjectGUID)
	}
}

func TestCalendarButton(t *testing.T) {
	b := Calendar("btn", "Pick a date", "DateGregorian", "", "1900", "2100")
	if b.ButtonCalendar == nil {
		t.Fatal("ButtonCalendar should not be nil")
	}
	if b.ButtonCalendar.Type != "DateGregorian" {
		t.Errorf("Type = %q", b.ButtonCalendar.Type)
	}
}

func TestNumberPickerButton(t *testing.T) {
	b := NumberPicker("btn", "Pick a number", "0", "100", "50")
	if b.ButtonNumberPicker == nil {
		t.Fatal("ButtonNumberPicker should not be nil")
	}
	if b.ButtonNumberPicker.MinValue != "0" {
		t.Errorf("MinValue = %q", b.ButtonNumberPicker.MinValue)
	}
}

func TestStringPickerButton(t *testing.T) {
	b := StringPicker("btn", "Pick", []string{"a", "b", "c"}, "a")
	if b.ButtonStringPicker == nil {
		t.Fatal("ButtonStringPicker should not be nil")
	}
	if len(b.ButtonStringPicker.Items) != 3 {
		t.Errorf("Items = %v", b.ButtonStringPicker.Items)
	}
}

func TestLocationButton(t *testing.T) {
	b := Location("btn", "Pick a location", "Picker", 35.6892, 51.3890)
	if b.ButtonLocation == nil {
		t.Fatal("ButtonLocation should not be nil")
	}
	if b.ButtonLocation.DefaultPointerLocation == nil {
		t.Fatal("DefaultPointerLocation should not be nil")
	}
}

func TestTextboxButton(t *testing.T) {
	b := Textbox("btn", "Enter text", "SingleLine", "String", "Type here...", "")
	if b.ButtonTextbox == nil {
		t.Fatal("ButtonTextbox should not be nil")
	}
	if b.ButtonTextbox.PlaceHolder != "Type here..." {
		t.Errorf("PlaceHolder = %q", b.ButtonTextbox.PlaceHolder)
	}
}

func TestSelectionButton(t *testing.T) {
	items := []models.ButtonSelectionItem{{Text: "Option 1"}, {Text: "Option 2"}}
	b := Selection("btn", "Pick", "sel_id", items, false, "1")
	if b.ButtonSelection == nil {
		t.Fatal("ButtonSelection should not be nil")
	}
	if len(b.ButtonSelection.Items) != 2 {
		t.Errorf("Items = %v", b.ButtonSelection.Items)
	}
}

func TestFtoa(t *testing.T) {
	if got := ftoa(0); got != "0" {
		t.Errorf("ftoa(0) = %q", got)
	}
	if got := ftoa(3.14); !endsWith(got, ".140000") {
		t.Errorf("ftoa(3.14) = %q", got)
	}
}

func TestStripAt(t *testing.T) {
	if got := stripAt("@user"); got != "user" {
		t.Errorf("stripAt(@user) = %q", got)
	}
	if got := stripAt("user"); got != "user" {
		t.Errorf("stripAt(user) = %q", got)
	}
}

func endsWith(s, suffix string) bool {
	if len(suffix) > len(s) {
		return false
	}
	return s[len(s)-len(suffix):] == suffix
}
