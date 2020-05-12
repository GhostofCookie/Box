package box

import (
	"testing"
)

func TestGetFolderInfo(t *testing.T) {
	t.Run("TestInvalidConfig", func(t *testing.T) {
		sdk := new(SDK)
		_, err := sdk.GetFolderInfo("0")
		if err != errConfig {
			t.Error("Expected config to be invalid")
		}
	})

	t.Run("TestValidConfig", func(t *testing.T) {
		sdk := setup()
		err := sdk.RequestAccessToken()
		if err != nil {
			t.Error("Expected config to have been set")
		}
		_, err = sdk.GetFolderInfo("0")
		if err != nil {
			t.Error("Expected to receive Folder info")
		}
	})
}
