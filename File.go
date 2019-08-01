package box

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"reflect"
	"strings"
	"time"
)

// GetFileInfo : Get information about a file.
func GetFileInfo(fileID string) (*FileObject, error) {
	RequestAccessToken()
	response, err := BoxRequest("GET", "https://api.box.com/2.0/files/"+fileID, nil, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fileObject := &FileObject{}
	json.Unmarshal(response, &fileObject)

	return fileObject, nil
}

// DownloadFile : Retrieves the actual data of the file. An optional version
// parameter can be set to download a previous version of the file.
func DownloadFile(fileID string, location string) error {
	RequestAccessToken()
	response, err := BoxRequest("GET", "https://api.box.com/2.0/files/"+fileID+"/content", nil, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	fInfo, err := GetFileInfo(fileID)
	file, err := os.Create(location + fInfo.Name)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	_, err = file.Write(response)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// UploadFile : Use the Upload API to allow users to add a new file. The user
// can then upload a file by specifying the destination folder for the file.
// If the user provides a file name that already exists in the destination
// folder, the user will receive an error.
func UploadFile(file interface{}, newName string, folderID string) (*PathCollection, error) {
	RequestAccessToken()

	t := reflect.TypeOf(file)

	var name string
	if t.Name() == "string" {
		name = file.(string)
	} else {
		name = newName
	}
	if newName == "" && name != "" {
		newName = name
	}

	var contents []byte
	if t.Name() == "" {
		contents = file.([]byte)
	} else {
		f, err := os.Open(name)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()

		contents, err = ioutil.ReadAll(f)
		if err != nil {
			log.Println(err)
		}
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		log.Println(err)
	}
	part.Write(contents)

	err = writer.WriteField("filename", name)
	if err != nil {
		log.Println(err)
	}

	err = writer.Close()
	if err != nil {
		log.Println(err)
	}

	headers := make(map[string]string)
	headers["Content-Type"] = writer.FormDataContentType()
	headers["Content-Length"] = string(body.Len())

	response, err := BoxRequest("POST",
		"https://upload.box.com/api/2.0/files/content?attributes={%22name%22:%22"+newName+"%22,%20%22parent%22:{%22id%22:%22"+folderID+"%22}}",
		body, headers)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	fileObject := &PathCollection{}
	json.Unmarshal(response, &fileObject)

	return fileObject, nil
}

// UploadFileVersion : Uploading a new file version is performed in the same
// way as uploading a file. This method is used to upload a new version of an
// existing file in a user’s account.
func UploadFileVersion(fileID string, newName string) {}

type Session struct {
	sessionID string
	fileSize  int32
}

func (s *Session) NewFile(folderID string, fileSize int32, fileName string) {}

func (s *Session) NewVersion(folderID string, fileSize int32, fileName string) {}

func (s *Session) UploadPart() {}

func (s *Session) ListParts(offset int, limit int) {}

func (s *Session) CommitUpload(partID string, offset int, size int32) {}

func (s *Session) Abort() {}

func PreflightCheck(name string, parentID string, size int32) bool {
	return true
}

// DeleteFile : Deletes a file in a specific folder with 'ID" fileID.
func DeleteFile(fileID string, etag string) error {
	RequestAccessToken()
	headers := make(map[string]string)
	headers["If-Match"] = etag
	_, err := BoxRequest("DELETE", "https://api.box.com/2.0/files/"+fileID, nil, headers)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// CopyFile : Copy a file. The version and a new name can be optionally supplied.
func CopyFile(fileID string, folderID, version string, name string) (*FileObject, error) {
	RequestAccessToken()
	bodyStr := `{"parent": {"id" : ` + folderID + `}`
	if version != "" {
		bodyStr += `, "version" : ` + version
	}
	if name != "" {
		bodyStr += `, "name" : ` + name
	}
	bodyStr += `}`
	body := strings.NewReader(bodyStr)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	response, err := BoxRequest("GET", "https://api.box.com/2.0/files/"+fileID+"/copy", body, headers)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fileObject := &FileObject{}
	json.Unmarshal(response, &fileObject)

	return fileObject, nil
}

// LockandUnlock : Sets a lock, which expires within a specific time.
func LockandUnlock(fileID string, expiresAt time.Time, preventDownload bool) {

}

// GetThumbnail : Get a thumbnail image for a file.
func GetThumbnail(fileID string, extension string, minHeight int32, minWidth int32) {

}

// GetEmbedLink : Returns information about the file with 'ID' fileID.
func GetEmbedLink(fileID string) (*EmbeddedFile, error) {
	RequestAccessToken()
	response, err := BoxRequest("GET", "https://api.box.com/2.0/files/"+fileID+"?fields=expiring_embed_link", nil, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fileObject := &EmbeddedFile{}
	json.Unmarshal(response, &fileObject)

	return fileObject, nil
}

func GetFileCollaborations() {

}

func GetFileComments() {

}

func GetFileTasks() {

}
