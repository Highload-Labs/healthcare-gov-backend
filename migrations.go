package healthcare_gov_backend

import "embed"

//go:embed migrations/*.sql
var MigrationFiles embed.FS
