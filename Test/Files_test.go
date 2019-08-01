package box

import (
	"reflect"
	"testing"

	box "github.com/GhostofCookie/Box"
)

func TestGetFileInfo(t *testing.T) {
	box.SetConfigFile("box_test_config.json")
	file, err := box.GetFileInfo("497817679326")
	if reflect.DeepEqual(file, &box.FileObject{}) || err != nil {
		t.Fail()
	}
}

func TestDownloadFile(t *testing.T) {
	box.SetConfigFile("box_test_config.json")
	err := box.DownloadFile("497817679326", "")
	if err != nil {
		t.Fail()
	}

}

func TestUploadFile(t *testing.T) {
	box.SetConfigFile("box_test_config.json")
	pathCollection, err := box.UploadFile("Files_test.go", "test.txt", "0")
	if reflect.DeepEqual(pathCollection, &box.PathCollection{}) || err != nil {
		t.Fail()
	}
}

func TestUploadFileVersion(t *testing.T) {
	box.SetConfigFile("box_test_config.json")

}

// Add Testing for chunk upload
/* Here */

func TestCopyFile(t *testing.T) {
	box.SetConfigFile("box_test_config.json")
	file, err := box.CopyFile("497817679326", "0", "", "test_copy.txt")
	if reflect.DeepEqual(file, &box.FileObject{}) || err != nil {
		t.Fail()
	}
}

func TestLockandUnlock(t *testing.T) {
	box.SetConfigFile("box_test_config.json")

}

func TestGetThumbnail(t *testing.T) {
	box.SetConfigFile("box_test_config.json")

}

func TestGetEmbedLink(t *testing.T) {
	box.SetConfigFile("box_test_config.json")

}

func TestGetFileCollaborations(t *testing.T) {
	box.SetConfigFile("box_test_config.json")

}

func TestGetFileComments(t *testing.T) {
	box.SetConfigFile("box_test_config.json")

}

func TestGetFileTasks(t *testing.T) {
	box.SetConfigFile("box_test_config.json")

}

func TestDeleteFile(t *testing.T) {
	box.SetConfigFile("box_test_config.json")

}
