package models

import "github.com/jinzhu/gorm"

var migrations []schemaMigrator

// SchemaMigrator alters the DB schema.
type schemaMigrator func(*gorm.DB)

// RegisterSchemaMigrator registers a SchemaMigrator to be called upon
// the first connection to the DB.
func registerSchemaMigrator(m schemaMigrator) {
	migrations = append(migrations, m)
}

// SchemaMigrate ensures db is migrated up to current version.
func SchemaMigrate(db *gorm.DB) {
	for _, mfn := range migrations {
		mfn(db)
	}
}
