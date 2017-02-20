package medialocker

import (
	xdg "github.com/casimir/xdg-go"
	"github.com/spf13/afero"
)

type FileSystem struct {
	afero.Fs
	xdg.App
}

var (
	defaultFileSystem FileSystem
	testFileSystem    FileSystem
	testFsMode        bool = false
)

// LocalFileSystem returns the root filesystem of the host
func LocalFileSystem() FileSystem {
	if testFsMode {
		return testFileSystem
	}

	return defaultFileSystem
}

func LocalFileExists(path string) bool {
	exist, err := afero.Exists(LocalFileSystem().Fs, path)
	if err != nil {
		return false
	}

	return exist
}

func init() {
	defaultFileSystem = FileSystem{
		Fs:  afero.NewOsFs(),
		App: xdg.App{Name: "medialocker"},
	}
}
