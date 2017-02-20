package medialocker

import (
	"github.com/Sirupsen/logrus"
	"github.com/BurntSushi/toml"
)

type Settings struct {
	DB dbSettings `toml:"Database"`
	Server serverSettings
	Logger loggerSettings
}

type dbSettings struct {
	DataPath string
	LogSQL, MemDB bool
}

type serverSettings struct {
	Bind string
	AutoBind, AutoCert bool
	CertPath, KeyPath string
}

func BuildSettingsTemplate(ctx AppContext) Settings {
	fs := ctx.FileSystem()

	return Settings{
		Logger: loggerSettings{
			LogLevel: logrus.InfoLevel,
			LogPath: fs.DataPath("locker.log"),
		},
		Server: serverSettings{
			AutoBind: true,
			AutoCert: true,
		},
		DB: dbSettings{
			DataPath: fs.DataPath("locker.db"),
		},
	}
}

func LoadSettings(ctx AppContext, settingsPath string) Settings {
	var settings Settings
	fs := ctx.FileSystem()
	log := ctx.Logger()

	if settingsPath == "" {
		settingsPath = fs.ConfigPath("locker.conf")
	}

	log.Must("Failed to ensure config file directory exist at %s.", settingsPath).Do(fs.EnsureFileDirectory(settingsPath))
	if fs.FileExists(settingsPath) {
		log.Must("Failed to load locker config file %s.", settingsPath).Do(toml.DecodeFile(settingsPath, &settings))
	} else {
		settings = BuildSettingsTemplate(ctx)
		settingsFile, err := fs.Create(settingsPath)
		log.Must("Failed to create new config file %s.", settingsPath).Do(err)
		defer settingsFile.Close()
		encoder := toml.NewEncoder(settingsFile)
		encoder.Encode(settings)
	}

	return settings
}
