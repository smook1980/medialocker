package medialocker

import "testing"

func TestSettings_LoadConfig(t *testing.T) {
	ctx := NewTestAppCtx()
	t.Run("writes the template when file does not exist", func(m *testing.T) {
		template := BuildSettingsTemplate(ctx)
		if result := LoadSettings(ctx, "/locker.conf"); result != template {
			m.Errorf("Expected result to match template settings. \nRESULT:\n%v\n\nEXPECTED:\n%v", result, template)
		}
	})
}
