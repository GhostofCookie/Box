package box

import (
	"testing"
)

func TestCreateDeleteFile(t *testing.T) {
	t.Run("TestValidConfig", func(t *testing.T) {
		sdk := setup()
		err := sdk.RequestAccessToken()
		if err != nil {
			t.Error("Expected config to have been set")
		}
		collection, err := sdk.UploadFile("File_test.go", "TestFile", "0")
		if err != nil {
			t.Error("Expected to receive Path Collection info")
		}

		if len(collection.Entries) > 0 {
			err = sdk.DeleteFile(collection.Entries[0].ID, "0")
			if err != nil {
				t.Error("Expected no error from delete")
			}
		}
	})
}
