package enums

// ReportReason is the reason code sent when reporting an object.
type ReportReason string

const (
	ReportReasonPorn     ReportReason = "Porn"
	ReportReasonViolence ReportReason = "Violence"
	ReportReasonSpam     ReportReason = "Spam"
	ReportReasonAbuse    ReportReason = "Abuse"
	ReportReasonOther    ReportReason = "Other"
)

// String implements fmt.Stringer.
func (r ReportReason) String() string { return string(r) }
