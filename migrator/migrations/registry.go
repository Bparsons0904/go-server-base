package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type Migration struct {
	ID          string
	Description string
	Migrate     func(tx *gorm.DB) error
	Rollback    func(tx *gorm.DB) error
}

var RegisteredMigrations []*gormigrate.Migration

func RegisterMigration(migration Migration) {
	gorMigration := &gormigrate.Migration{
		ID:       migration.ID,
		Migrate:  migration.Migrate,
		Rollback: migration.Rollback,
	}
	RegisteredMigrations = append(RegisteredMigrations, gorMigration)
}
