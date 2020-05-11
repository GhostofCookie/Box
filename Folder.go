package box

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

const folderURL = "https://api.box.com/2.0/folders/"

// CreateFolder creates a new folder under the parent folder that has 'ID' parentFolderID.
func (sdk *SDK) CreateFolder(name string, parentFolderID string) (*FolderObject, error) {
	body := strings.NewReader(`{"name":"` + name + `", "parent": {"id": "` + parentFolderID + `"}}`)

	response, err := sdk.Request("POST", folderURL, body, nil)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	folderObject := &FolderObject{}
	json.Unmarshal(response, &folderObject)

	return folderObject, nil
}

// GetFolderItems returns all the items contained inside the folder with 'ID' folderID.
func (sdk *SDK) GetFolderItems(folderID string, limit int, offset int) (*ItemCollection, error) {
	urlOpts := "/items?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset)
	response, err := sdk.Request("GET", folderURL+folderID+urlOpts, nil, nil)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	items := &ItemCollection{}
	json.Unmarshal(response, &items)

	return items, nil
}

// DeleteFolder deletes the folder who's 'ID' matches folderID.
func (sdk *SDK) DeleteFolder(folderID string) {
	_, err := sdk.Request("DELETE", folderURL+folderID+"?recursive=true", nil, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
