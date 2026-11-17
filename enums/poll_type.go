package enums

// PollType is the type of a poll (Regular vs Quiz).
type PollType string

const (
	PollTypeRegular PollType = "Regular"
	PollTypeQuiz    PollType = "Quiz"
)

// String implements fmt.Stringer.
func (p PollType) String() string { return string(p) }
