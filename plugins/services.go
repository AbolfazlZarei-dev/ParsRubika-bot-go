package plugins

// Names of the services published by the built-in plugins.
const (
	ServiceAccess = "access"
	ServiceUsers  = "users"
	ServiceKV     = "kv"
	ServiceI18n   = "i18n"
	ServiceFlags  = "flags"
	ServiceStats  = "analytics"
)

// AccessService is published by the access plugin. Other plugins use it
// to exempt administrators from limits and to guard admin commands.
type AccessService interface {
	IsAdmin(userID string) bool
	IsBlocked(userID string) bool
}

// UserDirectory is published by the users plugin.
type UserDirectory interface {
	// Chats returns every chat ID the bot has seen.
	Chats() []string
	// Count returns the number of known chats.
	Count() int
}

// IsAdmin asks the access service whether userID is an administrator.
// Without an access plugin nobody is an administrator.
func IsAdmin(m *Manager, userID string) bool {
	if a, ok := Service[AccessService](m, ServiceAccess); ok {
		return a.IsAdmin(userID)
	}
	return false
}
