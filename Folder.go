package box

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

// CreateFolder : Creates a new folder under the parent folder that has 'ID' parentFolderID.
func CreateFolder(name string, parentFolderID string) (*FolderObject, error) {
	RequestAccessToken()
	body := strings.NewReader(`{"name":"` + name + `", "parent": {"id": "` + parentFolderID + `"}}`)

	response, err := BoxRequest("POST", "https://api.box.com/2.0/folders", body, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	folderObject := &FolderObject{}
	json.Unmarshal(response, &folderObject)

	return folderObject, nil
}

// GetFolderItems : Returns all the items contained inside the folder with 'ID' folderID.
func GetFolderItems(folderID string, limit int, offset int) (*ItemCollection, error) {
	RequestAccessToken()

	response, err := BoxRequest("GET", "https://api.box.com/2.0/folders/"+folderID+"/items?limit="+strconv.Itoa(limit)+"&offset="+strconv.Itoa(offset), nil, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	items := &ItemCollection{}
	json.Unmarshal(response, &items)

	return items, nil
}

// DeleteFolder : Deletes the folder with 'ID' folderID.
func DeleteFolder(folderID string) {
	RequestAccessToken()
	_, err := BoxRequest("DELETE", "https://api.box.com/2.0/folders/"+folderID+"?recursive=true", nil, nil)
	if err != nil {
		log.Println(err)
		return
	}
}
