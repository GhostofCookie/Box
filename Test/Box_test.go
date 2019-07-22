package box

import (
	"reflect"
	"testing"

	box "github.com/GhostofCookie/Box"
)

func TestGetFileInfo(t *testing.T) {
	file, err := box.GetFileInfo("0")

	if reflect.DeepEqual(file, &box.FileObject{}) || err != nil {
		t.Fail()
	}
}
