package box

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
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
	response, err := sdk.request("GET", fileURL+fileID, nil, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fileObject := &FileObject{}
	json.Unmarshal(response, &fileObject)

	return fileObject, nil
}

// GetThumbnail gets a thumbnail image for the requested file.
func (sdk *SDK) GetThumbnail(fileID, extension string, minHeight, minWidth int) (image.Image, error) {
	opts := "?min_height=" + strconv.Itoa(minHeight) + "&min_width=" + strconv.Itoa(minWidth)
	response, err := sdk.request("GET", fileURL+fileID+"/thumbnail."+extension+opts, nil, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var image []byte
	json.Unmarshal(response, &image)

	switch extension {
	case "png":
		return png.Decode(bytes.NewReader(image))
	case "jpg":
	case "jpeg":
		return jpeg.Decode(bytes.NewReader(image))
	}
	return nil, nil
}

// CopyFile copies a file. The version and a new name can be optionally supplied.
func (sdk *SDK) CopyFile(fileID, folderID, name, version string) (*FileObject, error) {
	body := map[string]interface{}{"parent": map[string]string{"id": folderID}}
	if name != "" {
		body["name"] = name
	}
	if version != "" {
		body["version"] = version
	}
	payload, err := json.Marshal(body)

	headers := map[string]string{"Content-Type": "application/json"}
	response, err := sdk.request("POST", fileURL+fileID+"/copy", bytes.NewBufferString(string(payload)), headers)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fileObject := &FileObject{}
	json.Unmarshal(response, &fileObject)

	return fileObject, nil
}

// GetEmbedLink returns information about the file with 'ID' fileID.
func (sdk *SDK) GetEmbedLink(fileID string) (*EmbeddedFile, error) {
	sdk.RequestAccessToken()
	response, err := sdk.request("GET", fileURL+fileID+"?fields=expiring_embed_link", nil, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fileObject := &EmbeddedFile{}
	json.Unmarshal(response, &fileObject)

	return fileObject, nil
}

// UpdateFile updates the file with a new name and content.
func (sdk *SDK) UpdateFile(fileID, newFilename string) error {
	return errors.New("No implementation for (SDK)UpdateFile()")
}

// DeleteFile deletes a file in a specific folder with an 'ID' matching fileID.
func (sdk *SDK) DeleteFile(fileID, etag string) error {
	headers := make(map[string]string)
	headers["If-Match"] = etag
	_, err := sdk.request("DELETE", fileURL+fileID, nil, headers)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// DownloadFile : Retrieves the actual data of the file. An optional version
// parameter can be set to download a previous version of the file.
func (sdk *SDK) DownloadFile(fileID, location string) error {
	response, err := sdk.request("GET", fileURL+fileID+"/content", nil, nil)
	if err != nil {
		log.Println(err)
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
func (sdk *SDK) UploadFile(inFile interface{}, newFilename, folderID string) (*PathCollection, error) {
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

	headers := map[string]string{
		"Content-Type":   writer.FormDataContentType(),
		"Content-Length": string(body.Len()),
	}

	response, err := sdk.request("POST", uploadURL, body, headers)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	pathCollection := &PathCollection{}
	json.Unmarshal(response, &pathCollection)

	return pathCollection, nil
}

///////////////////////////////////////////////////////////////////////////////
// FILE VERSION
///////////////////////////////////////////////////////////////////////////////

// UploadFileVersion : Uploading a new file version is performed in the same
// way as uploading a file. This method is used to upload a new version of an
// existing file in a user’s account.
func (sdk *SDK) UploadFileVersion(fileID, newName string) {}

///////////////////////////////////////////////////////////////////////////////
// UPLOAD (CUNKED)
///////////////////////////////////////////////////////////////////////////////

// Session TODO: Add definition
type Session struct {
	sessionID string
	fileSize  int32
}

// NewFile TODO: Add definition
func (s *Session) NewFile(folderID string, fileSize int32, fileName string) {}

// NewVersion TODO: Add definition
func (s *Session) NewVersion(folderID string, fileSize int32, fileName string) {}

// UploadPart TODO: Add definition
func (s *Session) UploadPart() {}

// ListParts TODO: Add definition
func (s *Session) ListParts(offset int, limit int) {}

// CommitUpload TODO: Add definition
func (s *Session) CommitUpload(partID string, offset int, size int32) {}

// Abort TODO: Add definition
func (s *Session) Abort() {}

// PreflightCheck TODO: Add definition
func PreflightCheck(name string, parentID string, size int32) bool {
	return true
}
