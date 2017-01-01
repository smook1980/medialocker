package medialocker

import "testing"

func TestAppInit(t *testing.T) {
	app, err := AppInit(true)
	if err != nil {
		t.Error(err)
	}
	if app == nil {
		t.Error("Expected non-nil value returned for app, got nil.")
	}

	return
}
