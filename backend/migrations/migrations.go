package migrations

import "embed"

// Files contains the SQL migrations used by the built-in migration runner.
//
//go:embed *.sql
var Files embed.FS
