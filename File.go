package box

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	fileURL   = "https://api.box.com/2.0/files/"
	uploadURL = "https://upload.box.com/api/2.0/files/content"
)

// GetFileInfo : Get information about a file.
func (sdk *SDK) GetFileInfo(fileID string) (*FileObject, error) {
	response, err := sdk.Request("GET", fileURL+fileID, nil, nil)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	fileObject := &FileObject{}
	json.Unmarshal(response, &fileObject)

	return fileObject, nil
}

// DownloadFile : Retrieves the actual data of the file. An optional version
// parameter can be set to download a previous version of the file.
func (sdk *SDK) DownloadFile(fileID string, location string) error {
	response, err := sdk.Request("GET", fileURL+fileID+"/content", nil, nil)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	if location == "." || location == ".." || !strings.HasSuffix(location, "/") {
		location += "/"
	}

	fInfo, err := sdk.GetFileInfo(fileID)
	file, err := os.Create(location + fInfo.Name)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	defer file.Close()

	_, err = file.Write(response)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

// UploadFile uses the Upload API to allow users to add a new file. The user
// can then upload a file by specifying the destination folder for the file.
// If the user provides a file name that already exists in the destination
// folder, the user will receive an error.
func (sdk *SDK) UploadFile(inFile interface{}, newFilename string, folderID string) (*PathCollection, error) {
	fileInputType := reflect.TypeOf(inFile)

	var filename string
	var contents []byte
	if fileInputType.Name() == "" {
		// Passed a file object.
		contents = inFile.([]byte)
	} else if fileInputType.Name() == "string" {
		// Passed a filename rather than a file.
		filename = inFile.(string)
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalln(err)
			return nil, err
		}
		defer file.Close()

		contents, err = ioutil.ReadAll(file)
		if err != nil {
			log.Fatalln(err)
			return nil, err
		}
	}

	// Did not specify a new file name, so keep anme the same.
	if newFilename == "" && filename != "" {
		newFilename = filename
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", newFilename)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	part.Write(contents)

	fields := map[string]string{
		"filename":   filename,
		"attributes": `{"name":"` + newFilename + `","parent":{"id":"` + folderID + `"}}`,
	}

	for field, value := range fields {
		err = writer.WriteField(field, value)
		if err != nil {
			log.Fatalln(err)
			return nil, err
		}
	}

	err = writer.Close()
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	headers := make(map[string]string)
	headers["Content-Type"] = writer.FormDataContentType()
	headers["Content-Length"] = string(body.Len())

	response, err := sdk.Request("POST", uploadURL, body, headers)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	pathCollection := &PathCollection{}
	json.Unmarshal(response, &pathCollection)

	return fileObject, nil
}

// UploadFileVersion : Uploading a new file version is performed in the same
// way as uploading a file. This method is used to upload a new version of an
// existing file in a user’s account.
func (sdk *SDK) UploadFileVersion(fileID string, newName string) {}

// DeleteFile : Deletes a file in a specific folder with 'ID" fileID.
func (sdk *SDK) DeleteFile(fileID string, etag string) error {
	headers := make(map[string]string)
	headers["If-Match"] = etag
	_, err := sdk.Request("DELETE", fileURL+fileID, nil, headers)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

// CopyFile : Copy a file. The version and a new name can be optionally supplied.
func (sdk *SDK) CopyFile(fileID string, folderID, version string, name string) (*FileObject, error) {
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
	response, err := sdk.Request("GET", fileURL+fileID+"/copy", body, headers)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	fileObject := &FileObject{}
	json.Unmarshal(response, &fileObject)

	return fileObject, nil
}

// GetThumbnail gets a thumbnail image for the requested file.
func (sdk *SDK) GetThumbnail(fileID string, extension string, minHeight int, minWidth int) (interface{}, error) {
	opts := "?min_height=" + strconv.Itoa(minHeight) + "min_width=" + strconv.Itoa(minWidth)
	response, err := sdk.Request("GET", fileURL+fileID+"/thumbnail."+extension+opts, nil, nil)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	var image interface{}
	json.Unmarshal(response, &image)

	log.Println(image)

	return image, nil
}

// GetEmbedLink : Returns information about the file with 'ID' fileID.
func (sdk *SDK) GetEmbedLink(fileID string) (*EmbeddedFile, error) {
	sdk.RequestAccessToken()
	response, err := sdk.Request("GET", fileURL+fileID+"?fields=expiring_embed_link", nil, nil)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	fileObject := &EmbeddedFile{}
	json.Unmarshal(response, &fileObject)

	return fileObject, nil
}

///////////////////////////////////////////////////////////////////////////////
// CHUNK FILE
///////////////////////////////////////////////////////////////////////////////

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
