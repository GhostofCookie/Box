package box

import (
	"testing"
)

func TestSetConfigFile(t *testing.T) {
	sdk := new(SDK)
	sdk.NewConfigFromFile("config.json")
	sdk.RequestAccessToken()

	fileObj, _ := sdk.UploadFile("LICENSE", "TEST_FILE", "0")

	if len(fileObj.Entries) > 0 {
		sdk.DownloadFile(fileObj.Entries[0].ID, ".")
		sdk.DeleteFile(fileObj.Entries[0].ID, "0")
	}
}
