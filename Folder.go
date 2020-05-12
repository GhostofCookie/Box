package box

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

const folderURL = "https://api.box.com/2.0/folders/"

// GetFolderInfo gets the information for the requested folder ID
func (sdk *SDK) GetFolderInfo(folderID string) (*FolderObject, error) {
	response, err := sdk.request("GET", folderURL+folderID, nil, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	folder := &FolderObject{}
	json.Unmarshal(response, &folder)

	return folder, nil
}

// ListItemsInFolder returns all the items contained inside the folder with 'ID' folderID.
func (sdk *SDK) ListItemsInFolder(folderID string, limit int, offset int) (*ItemCollection, error) {
	urlOpts := "/items?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset)
	response, err := sdk.request("GET", folderURL+folderID+urlOpts, nil, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	items := &ItemCollection{}
	json.Unmarshal(response, &items)

	return items, nil
}

// CreateFolder creates a new folder under the parent folder that has 'ID' parentFolderID.
func (sdk *SDK) CreateFolder(name string, parentFolderID string) (*FolderObject, error) {
	body := strings.NewReader(`{"name":"` + name + `", "parent": {"id": "` + parentFolderID + `"}}`)

	response, err := sdk.request("POST", folderURL, body, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	folderObject := &FolderObject{}
	json.Unmarshal(response, &folderObject)

	return folderObject, nil
}

func (sdk *SDK) CopyFolder(folderID string, parentFolderID string, newName string) (*FolderObject, error) {
	body := map[string]interface{}{"parent": map[string]string{"id": parentFolderID}}
	if newName != "" {
		body["name"] = newName
	}
	payload, err := json.Marshal(body)

	headers := map[string]string{"Content-Type": "application/json"}
	response, err := sdk.request("POST", folderURL+folderID+"/copy", bytes.NewBufferString(string(payload)), headers)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	folderObject := &FolderObject{}
	json.Unmarshal(response, &folderObject)

	return folderObject, nil
}

func (sdk *SDK) UpdateFolder() {}

// DeleteFolder deletes the folder who's 'ID' matches folderID.
func (sdk *SDK) DeleteFolder(folderID string) {
	_, err := sdk.request("DELETE", folderURL+folderID+"?recursive=true", nil, nil)
	if err != nil {
		log.Println(err)
		return
	}
}
