package medialocker

import (
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/smook1980/medialocker/models"
)

const (
	DB_URI_TEMPLATE = "%s.db"
	// DB_URI_TEMPLATE = "file:%s.db?mode=rwc"
)

type DB interface {
	VideoRepository() interface{}
}

type DBConnection struct {
	*gorm.DB
}

func (db *DBConnection) Ping() error {
	return db.DB.DB().Ping()
}

func (db *DBConnection) VideoRepository() interface{} {
	return nil
}

type Closer func() error
type DBConnectionFactory func() (*DBConnection, error)

func NewDBConnectionFactory(log *Logger, c Config) (DBConnectionFactory, Closer) {
	var connect sync.Once
	var dbLock sync.Mutex

	var db *gorm.DB
	var dbUrl string
	logSQL := c.LogSQL

	if c.MemDB {
		dbUrl = ":memory:"
	} else {
		dbUrl = fmt.Sprintf(DB_URI_TEMPLATE, c.DbPath)
	}

	logger := log.WithField("db", dbUrl)
	logger = logger.WithField("module", "DBConnectionFactory")
	logger.Debug("Initializing DB Connection Factory...")

	closer := func() error {
		var close sync.Once
		var err error

		close.Do(func() {
			dbLock.Lock()
			defer dbLock.Unlock()
			if db != nil {
				err = db.Close()
			}

			logger.Debug("DB connection closed!")
		})

		return err
	}

	factory := func() (*DBConnection, error) {
		var err error

		connect.Do(func() {
			dbLock.Lock()
			defer dbLock.Unlock()
			logger.Debug("Connecting db...")
			db, err = gorm.Open("sqlite3", dbUrl)
			db.LogMode(logSQL)
			logger.Debug("Migrating Schema...")
			models.SchemaMigrate(db)
			logger.Debug("DB Opened...")
		})

		if err != nil {
			logger.Panicf("Failed to establish DB connection, this is a fatal error! %s", err)
		}

		dbLock.Lock()
		defer dbLock.Unlock()

		err = db.DB().Ping()
		if err != nil {
			logger.Debugf("Unable to ping connection, is it dead?  Attempting to reopen! %s", err)
			db.Close()
			db, err = gorm.Open("sqlite3", fmt.Sprintf(DB_URI_TEMPLATE, c.DbPath))

			if err != nil {
				logger.Panicf("Failed to establish DB connection, this is a fatal error! %s", err)
			}
		}

		return &DBConnection{DB: db}, nil
	}

	return factory, closer
}
