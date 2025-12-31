package migrations

import "embed"

// MigrationFS 嵌入的迁移文件系统
//go:embed *.sql
var MigrationFS embed.FS