package store

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/smook1980/medialocker/app"
)

var db *gorm.DB

func DB(ctx app.Context) (*gorm.DB, error) {
	var err error

	fmt.Printf("Connecting to a %s server at %s...", ctx.DbType(), ctx.DbURI())

	if db == nil {
		db, err = gorm.Open(ctx.DbType(), ctx.DbURI())
	}

	return db, err
}
