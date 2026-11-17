# `plugins` — Services

_Source: [`services.go`](../../ParsRubika/plugins/services.go)_

## Symbols

### `ServiceAccess`

```go
const ServiceAccess = "access"
```

Names of the services published by the built-in plugins.

### `ServiceUsers`

```go
const ServiceUsers = "users"
```

Names of the services published by the built-in plugins.

### `ServiceKV`

```go
const ServiceKV = "kv"
```

Names of the services published by the built-in plugins.

### `ServiceI18n`

```go
const ServiceI18n = "i18n"
```

Names of the services published by the built-in plugins.

### `ServiceFlags`

```go
const ServiceFlags = "flags"
```

Names of the services published by the built-in plugins.

### `ServiceStats`

```go
const ServiceStats = "analytics"
```

Names of the services published by the built-in plugins.

### `AccessService`

```go
type AccessService interface
```

AccessService is published by the access plugin. Other plugins use it
to exempt administrators from limits and to guard admin commands.

### `UserDirectory`

```go
type UserDirectory interface
```

UserDirectory is published by the users plugin.

### `IsAdmin`

```go
func IsAdmin(m *Manager, userID string) bool
```

IsAdmin asks the access service whether userID is an administrator.
Without an access plugin nobody is an administrator.
