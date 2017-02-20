package medialocker

import (
	"testing"

	"github.com/casimir/xdg-go"
	"github.com/spf13/afero"
)

func init() {
	testFileSystem = FileSystem{App: xdg.App{Name: "medialocker"}, Fs: afero.NewMemMapFs()}
}

func GetTestFileSystem() FileSystem {
	return testFileSystem
}

func ResetTestFileSystem() FileSystem {
	testFileSystem = FileSystem{App: xdg.App{Name: "medialocker"}, Fs: afero.NewMemMapFs()}

	return testFileSystem
}

func EnableTestFileSystem() {
	testFsMode = true
}

func DisableTestFileSystem() {
	testFsMode = false
}

func TestFileSystem(t *testing.T) {
	fs := GetTestFileSystem()
	if fs != testFileSystem {
		t.Error("GetTestFileSystem returned nil, expected test filesystem.")
	}
}

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	if ret == 0 {
		teardown()
	}
}

func setup() {
	EnableTestFileSystem()
	ResetTestFileSystem()
}

func teardown() {
	DisableTestFileSystem()
}
