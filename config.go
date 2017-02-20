package medialocker

import (
	"gopkg.in/ini.v1"
	"os"
	"path"
)

type Config struct {
	DbPath       string
	MemDB        bool
	LogSQL       bool
	ConfigPath   string
	LogPath      string
	Bind         string
	DebugLogging bool
}

var (
	EmptyConfig = Config{}
)

type Configuration func(*Config) error

func FileConfiguration(configPath string) Configuration {
	if configPath == "" {
		configPath = LocalFileSystem().ConfigPath("locker.conf")
	}

	if LocalFileExists(configPath) {
		return func(c *Config) error {
			cfg, err := ini.Load(configPath)

			if err != nil {
				return err
			}

			cfg.MapTo(c)

			return nil
		}
	} else {
		return func(c *Config) error {
			fn := DefaultConfiguration()
			err := fn(c)
			if err != nil {
				return err
			}

			confDir := path.Dir(configPath)
			err = LocalFileSystem().MkdirAll(confDir, os.ModeDir|os.ModePerm)
			if err != nil {
				return err
			}

			cfgFile, err := LocalFileSystem().Create(configPath)
			if err != nil {
				return err
			}

			cfg := ini.Empty()
			cfg.ReflectFrom(c)

			_, err = cfg.WriteTo(cfgFile)
			return err
		}
	}

}

// DefaultConfiguration will set any unset options a default value
func DefaultConfiguration() Configuration {
	return func(c *Config) error {
		if c.Bind == EmptyConfig.Bind {
			c.Bind = ":3000"
		}

		if c.ConfigPath == EmptyConfig.ConfigPath {
			c.ConfigPath = LocalFileSystem().ConfigPath("locker.conf")
		}

		if c.LogPath == EmptyConfig.LogPath {
			c.LogPath = LocalFileSystem().DataPath("medialocker.log")
		}

		if c.DbPath == EmptyConfig.DbPath {
			c.DbPath = LocalFileSystem().DataPath("db")
		}

		return nil
	}
}

func BuildConfig(opts ...Configuration) (Config, []error) {
	c := &Config{}
	errors := make([]error, 0)

	opts = append(opts, DefaultConfiguration())

	for _, fn := range opts {
		err := fn(c)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return *c, errors
}

func (c Config) withConfiguration(fn Configuration) (*Config, error) {
	err := fn(&c)

	return &c, err
}
