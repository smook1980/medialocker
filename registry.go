package medialocker

import (
	"github.com/Sirupsen/logrus"
)

type Registry struct {
	log              *logrus.Entry
	dbFactory        DBConnectionFactory
	dataStoreFactory DataStoreFactory
	closers          []Closer
}

func (r *Registry) DB() (*DBConnection, error) {
	return r.dbFactory()
}

func (r *Registry) Shutdown() {
	for _, fn := range r.closers {
		if err := fn(); err != nil {
			r.log.Errorf("Registry unable to clean shutdown! %s", err)
		}
	}
}

// TODO: Do not return nil for interface, return concrete type?
func (r *Registry) DataStore() (*DataStore, error) {
	return r.dataStoreFactory()
}

func NewRegistry(log *Logger, c Config) *Registry {
	registryLog := log.WithField("module", "registry")
	registryLog.Debug("Initializing service registry...")
	var closers []Closer
	db, closer := NewDBConnectionFactory(log, c)
	closers = append(closers, closer)

	return &Registry{
		log:       registryLog,
		dbFactory: db,
		closers:   closers,
	}
}
